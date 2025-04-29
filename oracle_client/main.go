package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/godror/godror" // Oracle驱动
)

type OracleClient struct {
	db *sql.DB
}

func main() {
	client := &OracleClient{}
	client.login()

	for {
		fmt.Println("=================================================================================================\n")
		client.executeQuery()
	}
}

func (c *OracleClient) login() {
	var (
		username string
		password string
		host     string
		port     string
		service  string
	)

	fmt.Print("请输入用户名: ")
	fmt.Scanln(&username)

	fmt.Print("请输入密码: ")
	fmt.Scanln(&password)

	fmt.Print("请输入主机地址: ")
	fmt.Scanln(&host)

	fmt.Print("请输入端口号(默认1521): ")
	fmt.Scanln(&port)
	if port == "" {
		port = "1521"
	}

	fmt.Print("请输入服务名: ")
	fmt.Scanln(&service)

	// 构建连接字符串
	connStr := fmt.Sprintf(`user="%s" password="%s" connectString="%s:%s/%s"`, username, password, host, port, service)

	// 连接数据库
	db, err := sql.Open("godror", connStr)
	if err != nil {
		log.Printf("连接数据库失败: %v\n", err)
		return
	}

	c.db = db
	fmt.Printf("登录oracle数据库成功, loginAddress[%s:%s], loginUser[%s], loginService[%s]\n", host, port, username, service)
}

func (c *OracleClient) executeQuery() {
	if c.db == nil {
		fmt.Println("请先登录数据库")
		return
	}

	fmt.Print("请输入SQL查询语句(不用带分号, 退出请输入exit): ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan() // 读取一行
	query := scanner.Text()
	if query == "exit" {
		c.logout()
		return
	}

	// 执行查询
	rows, err := c.db.Query(query)
	if err != nil {
		log.Printf("查询执行失败: %v\n", err)
		return
	}
	defer rows.Close()

	// 获取列名
	columns, err := rows.Columns()
	if err != nil {
		log.Printf("获取列名失败: %v\n", err)
		return
	}

	// 打印列名
	for _, col := range columns {
		fmt.Printf("%-20s", col)
	}
	fmt.Println()

	// 准备接收数据的切片
	values := make([]interface{}, len(columns))
	for i := range values {
		var v interface{}
		values[i] = &v
	}

	// 遍历结果集
	for rows.Next() {
		err := rows.Scan(values...)
		if err != nil {
			log.Printf("读取行数据失败: %v\n", err)
			continue
		}

		// 打印每一行数据
		for _, val := range values {
			v := *(val.(*interface{}))
			switch v := v.(type) {
			case []byte:
				fmt.Printf("%-20s", string(v))
			default:
				fmt.Printf("%-20v", v)
			}
		}
		fmt.Println()
	}

	if err = rows.Err(); err != nil {
		log.Printf("遍历结果集时出错: %v\n", err)
	}
}

func (c *OracleClient) logout() {
	if c.db != nil {
		err := c.db.Close()
		if err != nil {
			log.Printf("关闭数据库连接时出错: %v\n", err)
		} else {
			fmt.Println("已成功退出登录")
		}
		c.db = nil
	} else {
		fmt.Println("当前没有活动的数据库连接")
	}
	os.Exit(0)
}
