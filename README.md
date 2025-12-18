# Scraper-Project
It is a project written in the Go (Golang) language.

# Go Web Scraper & Screenshot Tool

Bu proje, Go (Golang) dili kullanılarak geliştirilmiş basit ve güvenilir bir **web scraper** uygulamasıdır.  
Uygulama, verilen bir URL’nin **ham HTML içeriğini** çekmekte ve aynı zamanda sayfanın **tam ekran görüntüsünü** alarak yerel diske kaydetmektedir.

---

## Özellikler

- Komut satırından URL alma
- Hedef siteye HTTP isteği gönderme
- Web sayfasının ham HTML içeriğini dosyaya kaydetme
- Tarayıcı tabanlı tam sayfa ekran görüntüsü alma
- Hata durumlarında kullanıcıyı bilgilendirme
- Basit, stabil ve genişletilebilir yapı

---

## Kullanılan Teknolojiler

- **Go (Golang)**
- **net/http** – HTTP istekleri
- **ios** – Dosya işlemleri
- **context** – Zaman aşımı ve kontrol
- **chromedp** – Tarayıcı otomasyonu ve ekran görüntüsü alma

---

## Kurulum

Öncelikle gerekli bağımlılıkları yükleyin:

```bash
go mod init scraper
go get github.com/chromedp/chromedp
