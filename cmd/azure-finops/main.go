package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Jessehoppus/azure-finops-dashboard-go/internal/adapters/azure/costquery"
)

func main() {
	var (
		scope     string
		fromStr   string
		toStr     string
		dimension string
		gran      string
	)

	flag.StringVar(&scope, "scope", "", "Escopo Azure (ex.: /subscriptions/<SUB_ID>)")
	flag.StringVar(&fromStr, "from", "", "Data inicial (YYYY-MM-DD)")
	flag.StringVar(&toStr, "to", "", "Data final (YYYY-MM-DD)")
	flag.StringVar(&dimension, "dimension", "ServiceName", "Dimensão para agrupamento (ServiceName, MeterCategory, ResourceGroup, TagKey:<tag>)")
	flag.StringVar(&gran, "granularity", "None", "Granularity: None|Daily|Monthly")
	flag.Parse()

	if scope == "" || fromStr == "" || toStr == "" {
		fmt.Println("Uso:")
		fmt.Println("  azure-finops trend --scope /subscriptions/<SUB_ID> --from 2025-09-24 --to 2025-10-24 --dimension ServiceName --granularity Monthly")
		os.Exit(2)
	}

	from, err := time.Parse("2006-01-02", fromStr)
	if err != nil {
		log.Fatalf("Data 'from' inválida: %v", err)
	}
	to, err := time.Parse("2006-01-02", toStr)
	if err != nil {
		log.Fatalf("Data 'to' inválida: %v", err)
	}

	ctx := context.Background()
	cli, err := costquery.NewClient(ctx)
	if err != nil {
		log.Fatalf("Erro de autenticação Azure: %v", err)
	}

	rows, headers, err := cli.CostByDimension(ctx, scope, from, to, dimension, gran)
	if err != nil {
		log.Fatalf("Erro na consulta: %v", err)
	}

	fmt.Printf("Período: %s .. %s | Escopo: %s\n", from.Format("2006-01-02"), to.Format("2006-01-02"), scope)
	for i, h := range headers {
		if i > 0 {
			fmt.Print("\t")
		}
		fmt.Print(h)
	}
	fmt.Println()

	for _, r := range rows {
		for i, c := range r {
			if i > 0 {
				fmt.Print("\t")
			}
			fmt.Print(c)
		}
		fmt.Println()
	}
}
