package reportpdf

import (
	"fmt"
	"os"

	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
)

func InitCollectPDF(theader []string, content [][]string, header string) {
	init := pdf.NewMaroto(consts.Portrait, consts.A4)
	init.SetPageMargins(20, 10, 20)

	buildHeading(init)
	init.Line(consts.MaxGridSum)
	buildTransactionList(init, theader, content, header)
	err := init.OutputFileAndClose("report/files/init.pdf")
	if err != nil {
		fmt.Println("Tidak Bisa Meyimpan PDF : ", err)
		os.Exit(1)
	}

	fmt.Println("Berhasil Mengenerate PDF")
}

func buildTransactionList(m pdf.Maroto, tableheadings []string, contents [][]string, header string) {

	m.Row(10, func() {
		m.Col(12, func() {
			m.Text(header, props.Text{
				Top:   2,
				Size:  13,
				Color: color.NewBlack(),
				Style: consts.Bold,
				Align: consts.Center,
			})
		})
	})

	m.TableList(tableheadings, contents, props.TableList{
		HeaderProp: props.TableListContent{
			Size: 9,
			// GridSizes: []uint{3, 7, 2},
		},
		ContentProp: props.TableListContent{
			Size: 8,
			// GridSizes: []uint{3, 7, 2},
		},
		Align:              consts.Left,
		HeaderContentSpace: 1,
		Line:               true,
	})
}

func buildHeading(m pdf.Maroto) {
	m.RegisterHeader(func() {
		m.Row(50, func() {
			m.Col(3, func() {
				_ = m.FileImage("images/pdf/logo_uty.png", props.Rect{
					Center:  true,
					Percent: 80,
				})
			})

			m.ColSpace(5)

			m.Col(4, func() {
				m.Text("Universitas Teknologi Yogyakarta", props.Text{
					Top:         10,
					Size:        10,
					Align:       consts.Left,
					Extrapolate: false,
					Style:       consts.Bold,
				})
				m.Text("Jl. Siliwangi Jl. Ring Road Utara, Jombor Lor, Sendangadi, Kec. Mlati, Kabupaten Sleman, Daerah Istimewa Yogyakarta", props.Text{
					Size:  8,
					Top:   14,
					Align: consts.Left,
					// Extrapolate: false,
				})
				m.Text("www.fundspoint.site", props.Text{
					Top:   30,
					Style: consts.BoldItalic,
					Size:  8,
					Align: consts.Left,
				})
			})

			m.TableList([]string{"Nama", "No Rekening", "Bank", "Nominal"}, [][]string{{"", "", "", ""}}, props.TableList{})
		})
	})
}
