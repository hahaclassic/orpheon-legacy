package tableoutput

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

type TablePrinter struct {
	style   table.Style
	headers []string
	rows    [][]any
}

type Option func(*TablePrinter)

func WithStyle(style table.Style) Option {
	return func(p *TablePrinter) {
		p.style = style
	}
}

func WithHeaders(headers []string) Option {
	return func(p *TablePrinter) {
		p.headers = headers
	}
}

func WithRows(rows [][]any) Option {
	return func(p *TablePrinter) {
		p.rows = rows
	}
}

func NewTablePrinter(options ...Option) *TablePrinter {
	tp := &TablePrinter{
		style: table.StyleDefault,
	}
	for _, option := range options {
		option(tp)
	}

	return tp
}

func (p *TablePrinter) SetHeaders(headers []string) {
	p.headers = headers
}

func (p *TablePrinter) SetRows(rows [][]any) {
	p.rows = rows
}

func (p *TablePrinter) PrintDefaultTable() {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	if len(p.headers) > 0 {
		tableHeader := make(table.Row, 0, len(p.headers))
		for i := range p.headers {
			tableHeader = append(tableHeader, p.headers[i])
		}
		t.AppendHeader(tableHeader)
	}

	if len(p.rows) > 0 {
		for _, row := range p.rows {
			t.AppendRow(row)
		}
	}

	t.SetStyle(p.style)
	t.Render()
}

func (p *TablePrinter) PrintTable(headers []string, rows [][]any) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	if len(headers) > 0 {
		tableHeader := make(table.Row, 0, len(headers))
		for i := range headers {
			tableHeader = append(tableHeader, headers[i])
		}
		t.AppendHeader(tableHeader)
	}

	if len(rows) > 0 {
		for _, row := range rows {
			t.AppendRow(row)
		}
	}

	t.SetStyle(p.style)
	t.Render()
}

func PrintTable(style table.Style, headers []string, rows [][]any) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	tableHeader := make(table.Row, 0, len(headers))
	for i := range headers {
		tableHeader = append(tableHeader, headers[i])
	}
	t.AppendHeader(tableHeader)

	for _, row := range rows {
		t.AppendRow(row)
	}
	t.SetStyle(style)
	t.Render()
}
