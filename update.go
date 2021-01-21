package mnc

import (
	"github.com/zofan/go-country"
	"github.com/zofan/go-fwrite"
	"github.com/zofan/go-req"
	"github.com/zofan/go-xmlre"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func Update() error {
	var (
		httpClient = req.New(req.DefaultConfig)
		list       []MNC
	)

	var (
		tableRe = xmlre.Compile(`(<h3>.*?)</table>`)
		nameRe  = xmlre.Compile(`<h3>(.*?)</`)
		rowRe   = xmlre.Compile(`<tr>(.*?)</tr>`)
		colsRe  = xmlre.Compile(`<td>(.*?)</td>`)
	)

	resp := httpClient.Get(`https://cellidfinder.com/mcc-mnc`)
	if resp.Error() != nil {
		return resp.Error()
	}
	body := string(resp.ReadAll())
	body = strings.ReplaceAll(body, `&nbsp;`, ` `)

	for _, table := range tableRe.FindAllStringSubmatch(body, -1) {
		cn := nameRe.FindStringSubmatch(table[1])

		for _, operator := range rowRe.FindAllStringSubmatch(table[1], -1) {
			cols := colsRe.FindAllStringSubmatch(operator[1], -1)

			if len(cols) == 0 {
				continue
			}

			if strings.TrimSpace(cols[4][1]) != statusActive {
				continue
			}

			mcc, _ := strconv.Atoi(strings.TrimSpace(cols[0][1]))
			mnc, _ := strconv.Atoi(strings.TrimSpace(cols[1][1]))

			cCode := ``
			if c := country.ByName(cn[1]); c != nil {
				cCode = c.Alpha3
			} else {
				continue
			}

			op := MNC{
				Country: cCode,

				Network:  strings.TrimSpace(cols[2][1]),
				Operator: strings.TrimSpace(cols[3][1]),

				MCC: mcc,
				MNC: mnc,
			}

			list = append(list, op)
		}
	}

	var tpl []string

	tpl = append(tpl, `package mnc`)
	tpl = append(tpl, ``)
	tpl = append(tpl, `var List = []MNC{`)

	for _, op := range list {
		tpl = append(tpl, `	{`)
		tpl = append(tpl, `		Country:  "`+op.Country+`",`)
		tpl = append(tpl, `		Network:  "`+op.Network+`",`)
		tpl = append(tpl, `		Operator: "`+op.Operator+`",`)
		tpl = append(tpl, `		MCC:      `+strconv.Itoa(op.MCC)+`,`)
		tpl = append(tpl, `		MNC:      `+strconv.Itoa(op.MNC)+`,`)
		tpl = append(tpl, `	},`)
	}

	tpl = append(tpl, `}`)

	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)

	return fwrite.WriteRaw(dir+`/mnc_db.go`, []byte(strings.Join(tpl, "\n")))
}
