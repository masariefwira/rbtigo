package email

import (
	"fmt"

	"github.com/ikalkali/rbti-go/entity/models"
)

func MapEmailKeterlambatan(input models.PeminjamanResponse) string {
	var (
		nama    = input.Nama
		tenggat = input.TenggatPengembalian.Format("January 2, 2006")
		denda   = input.Denda
		buku    = input.DetailPeminjaman
		table   string
	)

	head := `
		<html lang="en">
		<body style="font-family: sans-serif; width: 30%%; margin: auto; padding: 10px;">
			<div style="padding: 10px; background-image: linear-gradient( #039dfc, #4aa3db); color: white;">
				<h1>RBTI Universitas Brawijaya</h1>
			</div>
			<div style="padding: 10px">
			<h2>REMINDER</h2>
			<hr>
			<h3>Halo %v!</h3>
			<p>Ada beberapa buku nih yang harus kamu kembalikan ke perpustakaan</p>
			<table style="border-collapse: collapse;">
			<tr>
            <th style="border: 1px solid black; padding: 5px;">
                Judul
            </th>
            <th style="border: 1px solid black; padding: 5px;">
                Penulis
            </th>
        </tr>
	`

	headFormatted := fmt.Sprintf(head, nama)

	tableTemplate := `
		<tr>
		<td style="border: 1px solid black; padding: 5px;">
			%v
		</td>
		<td style="border: 1px solid black; padding: 5px;">
			%v
		</td>
		</tr>
	`

	foot := `
		</table>
			<p>Segera kembalikan agar denda tidak terus bertambah</p>
			<p style="color: red;" >Tenggat pengembalian %v</p>
			<p style="color: red;" >Total Denda <strong>Rp. %v</strong></p>
			<p style="color: red; font-size: 10px;" >Denda per hari <strong>Rp. 5000</strong></p>
			<p>Apabila ada kendala silahkan menghubungi 083473493</p>
			<h3>Terima Kasih!</h3>
		</div>
		</body>
		</html>
	`

	footFormatted := fmt.Sprintf(foot, tenggat, denda)

	for _, data := range buku {
		tableTemp := fmt.Sprintf(tableTemplate, data.Judul, data.Penerbit)
		table = table + tableTemp
	}

	return headFormatted + table + footFormatted
}

func MapEmailReminder(input models.PeminjamanResponse) string {
	var (
		nama    = input.Nama
		tenggat = input.TenggatPengembalian.Format("January 2, 2006")
		buku    = input.DetailPeminjaman
		table   string
	)

	head := `
		<html lang="en">
		<body style="font-family: sans-serif; width: 30%%; margin: auto; padding: 10px;">
			<div style="padding: 10px; background-image: linear-gradient( #039dfc, #4aa3db); color: white;">
				<h1>RBTI Universitas Brawijaya</h1>
			</div>
			<div style="padding: 10px">
			<h2>REMINDER</h2>
			<hr>
			<h3>Halo %v!</h3>
			<p>Segera kembalikan buku dibawah ini sebelum dikenakan denda</p>
			<table style="border-collapse: collapse;">
			<tr>
            <th style="border: 1px solid black; padding: 5px;">
                Judul
            </th>
            <th style="border: 1px solid black; padding: 5px;">
                Penulis
            </th>
        </tr>
	`

	headFormatted := fmt.Sprintf(head, nama)

	tableTemplate := `
		<tr>
		<td style="border: 1px solid black; padding: 5px;">
			%v
		</td>
		<td style="border: 1px solid black; padding: 5px;">
			%v
		</td>
		</tr>
	`

	foot := `
		</table>
			<p>Kembalikan sebelum %v agar tidak dikenakan denda</p>
			<p>Apabila ada kendala silahkan menghubungi 083473493</p>
			<h3>Terima Kasih!</h3>
		</div>
		</body>
		</html>
	`

	footFormatted := fmt.Sprintf(foot, tenggat)

	for _, data := range buku {
		tableTemp := fmt.Sprintf(tableTemplate, data.Judul, data.Penerbit)
		table = table + tableTemp
	}

	return headFormatted + table + footFormatted
}
