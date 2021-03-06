package qsc

import (
	"bytes"
	"fmt"
	"github.com/QOSGroup/kepler/cert"
	keplercmd "github.com/QOSGroup/kepler/cmd"
	bacc "github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/context"
	"github.com/QOSGroup/qbase/txs"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/account"
	"github.com/QOSGroup/qos/mapper"
	"github.com/QOSGroup/qos/types"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/crypto"
)

// create QSC
type TxCreateQSC struct {
	Creator     btypes.Address        `json:"creator"`     //QSC创建账户
	Extrate     string                `json:"extrate"`     //qcs:qos汇率(amino不支持binary形式的浮点数序列化，精度同qos erc20 [.0000])
	QSCCA       *cert.Certificate     `json:"qsc_crt"`     //CA信息
	Description string                `json:"description"` //描述信息
	Accounts    []*account.QOSAccount `json:"accounts"`
}

func (tx TxCreateQSC) ValidateData(ctx context.Context) error {
	// CA校验
	if tx.QSCCA == nil {
		return errors.New("QSCCA is empty")
	}
	subj, ok := tx.QSCCA.CSR.Subj.(cert.QSCSubject)
	if !ok {
		return errors.New("invalid QSCSubject")
	}
	if subj.ChainId != ctx.ChainID() {
		return errors.New(fmt.Sprintf("chainId %s not matches %s ", subj.ChainId, ctx.ChainID()))
	}
	baseMapper := ctx.Mapper(mapper.BaseMapperName).(*mapper.MainMapper)
	rootCA := baseMapper.GetRootCA()
	if !keplercmd.VerityCrt([]crypto.PubKey{rootCA}, *tx.QSCCA) {
		return errors.New("invalid CA")
	}

	// QSC不存在
	qscMapper := ctx.Mapper(QSCMapperName).(*QSCMapper)
	if qscMapper.Exists(subj.Name) {
		return errors.New(fmt.Sprintf("%s already exists", subj.Name))
	}

	// creator账户存在
	accountMapper := ctx.Mapper(bacc.AccountMapperName).(*bacc.AccountMapper)
	creator := accountMapper.GetAccount(tx.Creator)
	if nil == creator {
		return errors.New("Creator account not exists")
	}

	_, ok = creator.(*account.QOSAccount)
	if !ok {
		return errors.New("Creator account is not a QOSAccount")
	}

	// accounts校验
	for _, account := range tx.Accounts {
		if account.QOS.NilToZero().GT(btypes.ZeroInt()) {
			return errors.New(fmt.Sprintf("invalid Accounts, %s QOS must be zero", account.AccountAddress))
		}
		if len(account.QSCs) != 1 || account.QSCs[0].Name != subj.Name {
			return errors.New(fmt.Sprintf("invalid Accounts, %s len(QSCs) must be 1 and QSCs[0].Name must be %s", account.AccountAddress, subj.Name))
		}
		if !account.QSCs[0].Amount.NilToZero().GT(btypes.ZeroInt()) {
			return errors.New(fmt.Sprintf("invalid Accounts, %s QSCs[0].Amount must gt zero", account.AccountAddress))
		}
	}

	return nil
}

func (tx TxCreateQSC) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.ABCICodeOK,
	}

	qscInfo := types.NewQSCInfoWithQSCCA(tx.QSCCA)
	qscInfo.Extrate = tx.Extrate
	qscInfo.Description = tx.Description

	// 保存QSC
	qscMapper := ctx.Mapper(QSCMapperName).(*QSCMapper)
	qscMapper.SaveQsc(&qscInfo)

	// 保存账户信息
	accountMapper := ctx.Mapper(bacc.AccountMapperName).(*bacc.AccountMapper)
	if qscInfo.Banker != nil {
		banker := qscInfo.Banker
		if nil == accountMapper.GetAccount(banker) {
			accountMapper.SetAccount(accountMapper.NewAccountWithAddress(banker))
		}
	}
	for _, acc := range tx.Accounts {
		if a := accountMapper.GetAccount(acc.AccountAddress); a != nil {
			qosAccount := a.(*account.QOSAccount)
			qosAccount.QSCs = qosAccount.QSCs.Plus(acc.QSCs)
			accountMapper.SetAccount(qosAccount)
		} else {
			accountMapper.SetAccount(acc)
		}
	}

	return
}

func (tx TxCreateQSC) GetSigner() []btypes.Address {
	return []btypes.Address{tx.Creator}
}

func (tx TxCreateQSC) CalcGas() btypes.BigInt {
	return btypes.ZeroInt()
}

func (tx TxCreateQSC) GetGasPayer() btypes.Address {
	return tx.Creator
}

func (tx TxCreateQSC) GetSignData() (ret []byte) {
	ret = append(ret, tx.Creator...)
	ret = append(ret, tx.Extrate...)
	ret = append(ret, cdc.MustMarshalBinaryBare(tx.QSCCA)...)
	ret = append(ret, tx.Description...)

	for _, account := range tx.Accounts {
		ret = append(ret, fmt.Sprint(account)...)
	}

	return
}

// issue QSC
type TxIssueQSC struct {
	QSCName string         `json:"qsc_name"` //币名
	Amount  btypes.BigInt  `json:"amount"`   //金额
	Banker  btypes.Address `json:"banker"`   //banker地址
}

func (tx TxIssueQSC) ValidateData(ctx context.Context) error {
	// QscName不能为空
	if len(tx.QSCName) < 0 {
		return errors.New("QSCName is empty")
	}

	// Amount大于0
	if !tx.Amount.GT(btypes.ZeroInt()) {
		return errors.New("Amount is lte zero")
	}

	// QSC存在
	qscMapper := ctx.Mapper(QSCMapperName).(*QSCMapper)
	qscInfo := qscMapper.GetQsc(tx.QSCName)
	if nil == qscInfo {
		return errors.New(fmt.Sprintf("QSCInfo of %s not exists", tx.QSCName))
	}

	// QSC名称一致
	if tx.QSCName != qscInfo.Name {
		return errors.New("wrong QSCName")
	}

	// qscInfo banker存在
	if qscInfo.Banker == nil {
		return errors.New("Banker not exists")
	}

	// banker 地址一致
	if !bytes.Equal(tx.Banker, qscInfo.Banker) {
		return errors.New("wrong Banker address")
	}

	return nil
}

func (tx TxIssueQSC) Exec(ctx context.Context) (result btypes.Result, crossTxQcp *txs.TxQcp) {
	result = btypes.Result{
		Code: btypes.ABCICodeOK,
	}

	accountMapper := ctx.Mapper(bacc.AccountMapperName).(*bacc.AccountMapper)

	banker := accountMapper.GetAccount(tx.Banker).(*account.QOSAccount)
	banker.QSCs = banker.QSCs.Plus(types.QSCs{btypes.NewBaseCoin(tx.QSCName, tx.Amount)})
	accountMapper.SetAccount(banker)

	return
}

func (tx TxIssueQSC) GetSigner() []btypes.Address {
	return []btypes.Address{tx.Banker}
}

func (tx TxIssueQSC) CalcGas() btypes.BigInt {
	return btypes.ZeroInt()
}

func (tx TxIssueQSC) GetGasPayer() btypes.Address {
	return tx.Banker
}

func (tx TxIssueQSC) GetSignData() (ret []byte) {
	ret = append(ret, tx.QSCName...)
	ret = append(ret, btypes.Int2Byte(tx.Amount.Int64())...)
	ret = append(ret, tx.Banker...)

	return
}
