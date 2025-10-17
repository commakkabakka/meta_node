package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

/*
	题目2：事务语句
		假设有两个表： accounts 表（包含字段 id 主键， balance 账户余额）和 transactions 表（包含字段 id 主键， from_account_id 转出账户ID， to_account_id 转入账户ID， amount 转账金额）。
	要求 ：
		编写一个事务，实现从账户 A 向账户 B 转账 100 元的操作。在事务中，需要先检查账户 A 的余额是否足够，如果足够则从账户 A 扣除 100 元，向账户 B 增加 100 元，并在 transactions 表中记录该笔转账信息。如果余额不足，则回滚事务。
*/

type Account struct {
	gorm.Model
	Balance float64
}

type Transaction struct {
	gorm.Model
	FromAccountID uint
	ToAccountID   uint
	Amount        float64
}

func Test32() {
	db, err := ConnectMysql()
	if err != nil {
		fmt.Println(err)
		return
	}

	db.AutoMigrate(&Account{})
	db.AutoMigrate(&Transaction{})

	//var accounts []Account = []Account{
	//	{Balance: 500},
	//	{Balance: 400},
	//	{Balance: 300},
	//	{Balance: 200},
	//	{Balance: 100},
	//}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	//err = gorm.G[[]Account](db).Create(ctx, &accounts)
	//if err != nil {
	//	fmt.Println(err)
	//}

	db.Transaction(func(tx *gorm.DB) error {

		accA, err := gorm.G[Account](db).Where("id = ?", 1).First(ctx)
		if err != nil {
			fmt.Println(err)
			return err
		}

		if accA.Balance < 100 {
			fmt.Println("余额不足")
			return errors.New("余额不足")
		}

		accB, err := gorm.G[Account](db).Where("id = ?", 2).First(ctx)
		if err != nil {
			fmt.Println(err)
			return err
		}

		//accA.Balance = accA.Balance - 100
		//accB.Balance = accB.Balance + 100

		// 零值不会更新 - 这里会有问题
		//_, err = gorm.G[Account](db).Where("id = ?", 1).Updates(ctx, accA)
		_, err = gorm.G[Account](db).Where("id = ?", 1).Update(ctx, "Balance", accA.Balance-100)
		if err != nil {
			fmt.Println(err)
			return err
		}

		// 零值不会更新 - 这里会有问题
		//_, err = gorm.G[Account](db).Where("id = ?", 2).Updates(ctx, accB)
		_, err = gorm.G[Account](db).Where("id = ?", 2).Update(ctx, "Balance", accB.Balance+100)
		if err != nil {
			fmt.Println(err)
			return err
		}

		var transaction Transaction = Transaction{
			FromAccountID: accA.ID,
			ToAccountID:   accB.ID,
			Amount:        100,
		}
		err = gorm.G[Transaction](db).Create(ctx, &transaction)
		if err != nil {
			fmt.Println(err)
			return err
		}

		return nil
	})
}
