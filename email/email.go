package email

import (
	"net/smtp"
	"strings"
)

var MESSAGE = `
<html lang="en">
<head>
	<meta charset="UTF-8" />
	<meta http-equiv="X-UA-Compatible" content="IE=edge" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
	<title>Document</title>
</head>
<body style="text-align: center; font-family: sans-serif">
	<h1>
		DAPATKAN PRIA BERBULU DADA <br />
		DI SEKITAR ANDA SEKARANG!
	</h1>
	<img
		src="https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTznL-iMLmbyqXdga9uXMD1Buu-GluWdKV6XA&usqp=CAU"
		alt=""
		style="width: 300px"
	/>
	<h3>Ingin pria berbulu dada untuk menemani anda seharian?</h3>
	<h4>
		Tidak usah pusing lagi! Kami menyediakan jasa escort pria bulu dada!
	</h4>
	<ul style="list-style-type: none">
		<li>Ongkos dimulai dari Rp. 100.000 per jam</li>
		<li>Tebal bulu dada diatas 5cm, <strong>LEBAT</strong></li>
		<li>Senang diajak segalanya ( ͡° ͜ʖ ͡°)</li>
	</ul>
	<h3>TUNGGU APA LAGI?</h3>
	<button
		style="
			background-color: yellow;
			padding: 10px;
			border: none;
			border: 1px solid black;
			cursor: pointer;
		"
	>
		PESAN SEKARANG
	</button>
</body>
</html>
`

var recipients = "ikal.ikhwan@gmail.com;sefira.icl@gmail.com;annis.icl@gmail.com;dianakamilia.icl@gmail.com;farahaprilita.icl@gmail.com;nadiaimee.icl@gmail.com;novaputera.icl@gmail.com;aliyaqurrota.icl@gmail.com;febrio.akbar5@gmail.com"

func SendToMail(to, subject, body string) error {
	var (
		user     = "rbti@ikal.app"
		password = "cYqH98L3cR8yQX8"
		host     = "smtpdm-ap-southeast-1.aliyun.com:80"
		mailtype = "html"
	)

	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var content_type string
	if mailtype == "html" {
		content_type = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
	}
	msg := []byte("To: " + to + "\r\nFrom: RBTI Universitas Brawijaya " + user + "\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	send_to := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, send_to, msg)
	return err
}

// func main() {
// 	user := "rbti@ikal.app"
// 	password := "cYqH98L3cR8yQX8"
// 	host := "smtpdm-ap-southeast-1.aliyun.com:80"
// 	to := recipients
// 	subject := "Kesempatan Yang Sangat Menarik!"
// 	body := MESSAGE
// 	fmt.Println("send email")
// 	err := SendToMail(to, subject, body, "html")
// 	if err != nil {
// 		fmt.Println("Send mail error!")
// 		fmt.Println(err)
// 	} else {
// 		fmt.Println("Send mail success!")
// 	}
// }
