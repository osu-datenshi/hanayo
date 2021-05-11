# Hanayo

Original source dari [Ripple Hanayo](https://github.com/osuripple/hanayo)

Source ini digunakan untuk bagian website nya

# Build
Untuk Build sekarang sudah support untuk latest go version (1.16.x)
1. Clone terlebih dahulu menggunakan `git clone https://github.com/osu-datenshi/hanayo` 
2. Masuk ke foldernya `cd hanayo`
3. Update submodules menggunakan `git submodule init && git submodule update`
4. Setelah itu build menggunakan `go build` tunggu nanti ada muncul downloading packages hingga selesai.
5. Hasil build bisa berbeda, jika kalian building menggunakan windows nama file akan menjadi `hanayo.exe` jika linux menjadi `hanayo` saja
6. Jalankan filenya menggunakan terminal only! contoh di linux `./hanayo` ikuti petunjuknya ketik `I agree` lalu selesai!

Beberapa hal yang harus dibuat sesudah menjalankannya :
- Buatlah folder `profbackgrounds` di dalam folder `static`
- Folder submodules `static` direpository diupload ke hostingan berbeda, karena isinya full file-static dan harus dipisahkan sebagai CDN (Content Delivery Network) jika punya server CDN sendiri, jika ga punya tidak perlu diupload, silahkan test di local environment masing-masing

Beberapa config yang harus diperhatikan adalah :
- DSN = databasenya (format : `user:password@tcp(localhost:3306)/database`)
- API Secret = key untuk authorized [API Backend](https://github.com/osu-datenshi/api)
- PASTIKAN EDIT DOMAIN YANG ADA DI SELURUH TEMPLATES (commit [e36e900](https://github.com/osu-datenshi/hanayo/commit/acd44a52ce6df3228984ea5ccd41c4b155ac31e1))
- PASTIKAN GOOGLE CAPTCHA AKTIF (BIAR NO SPAM)

### kontribusi

buat yg ingin serius kontribusi silahkan pm troke buat di invite ke orgz