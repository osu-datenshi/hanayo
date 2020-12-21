# Hanayo

Original source dari [Ripple Hanayo](https://github.com/osuripple/hanayo) dan sedikit modifikasi dari osu!thailand source

Source ini digunakan untuk bagian website nya

cara build `go get https://github.com/osu-datenshi/hanayo` dan setelah itu `go build`

Reminder : Static folder dipisahkan agar mempermudah proses development, dan untuk folder profile background (`/static/profbackgrounds`) perlu dibuatkan manual di dalam vps, sedangkan folder utama di upload ke server yg lain agar website load lebih cepat seperti CDN, dan folder `website-docs` juga ikut dipisahkan

hal yang harus diperhatikan adalah :

- DSN = databasenya (format : `user:password@tcp(localhost:3306)/database`)
- Mailgun = harus punya akun mailgun (daftar gratis tapi harus pake CC) ini fungsi buat reset password nanti dikirim ke email
- API Secret = api secret nya buat api backend
- PASTIKAN EDIT DOMAIN YANG ADA DI SELURUH TEMPLATES (commit [e36e900](https://github.com/osu-datenshi/hanayo/commit/acd44a52ce6df3228984ea5ccd41c4b155ac31e1))
- PASTIKAN GOOGLE CAPTCHA AKTIF (BIAR NO SPAM)

### kontribusi

buat yg ingin serius kontribusi silahkan pm troke buat di invite ke orgz