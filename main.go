package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/chromedp/chromedp" // İstediğimiz sitenin ekran görüntüsünü almak için githubdan
	// bu kütüphaneyi kullanıyoruz çünkü go nun kendi kütüphanelerinde bu özellik yoktur
)

// Bu programın amacı şu şekilde :
// Kullanıcıdan bir url istiyoruz ve girdiği urlnin html ve ekran görüntüsünü dosyalara kaydediyoruz
// HTML bilgisini HtmlBilgisi.html adlı dosyaya kaydediyoruz
// Ekran görüntüsünü ise EkranGoruntusu.png ye kaydediyoruz
// Bu dosyalarda projemizin bulunduğu dosya klasörüne kaydediliyor

func main() {
	// CLI parametreleri ile URL, timeout ve verbose alma
	url := flag.String("url", "", "Hedef URL (zorunlu)")
	timeout := flag.Int("t", 30, "Timeout süresi saniye cinsinden")
	verbose := flag.Bool("v", false, "Detaylı log modu")

	flag.Parse()

	// URL boşsa hata verip çık
	if *url == "" {
		fmt.Println("Kullanim : go run main.go -url <hedef URL> [-t <timeout>] [-v]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Hedef url yi terminalde belirtiyoruz
	fmt.Println("Hedef URL :", *url)

	// Eğer verbose aktifse bilgileri göster
	if *verbose {
		fmt.Println("Hedef URL :", *url)
		fmt.Println("Timeout :", *timeout, "saniye")
	}

	// Evet şimdide hedef url mize istek attıktan sonrada HTML bilgisini çekme kısmına çıkıyoruz

	// Bu kod verilen adrese bir GET isteği atıyor ve eğer hata olursa ekrana yazıp işlemi bitiriyor
	// hata yoksa cevabı alıyor ve bağlantıyı düzgün şekilde kapatıyor. Yani kısaca, bir siteye bağlanıp içeriğini okumaya hazırlık yapıyor.
	resp, err := http.Get(*url)
	if err != nil {
		fmt.Println("Siteye baglanilamadi :", err)
		return
	}
	defer resp.Body.Close()

	// Bu kodun yaptığı şey de şu sunucudan gelen cevabın durum kodunu kontrol ediyor. Eğer durum kodu 200 OK değilse
	// yani istek başarılı olmamışsa, ekrana "HTTP hata kodu alındı" mesajını ve gelen hatanın durumunu yazdırıyor, ardından işlemi sonlandırıyor.
	// yani bu parça, bağlantı kurulsa bile sayfanın doğru şekilde dönüp dönmediğini kontrol eden bir güvenlik adımı.
	if resp.StatusCode != http.StatusOK {
		fmt.Println("HTTP hata kodu alindi :", resp.Status)
		return
	}

	// bu kodun yaptığı şey de şu sunucudan gelen cevabın gövdesini tamamen okuyup body değişkenine atıyor. Eğer okuma sırasında bir hata
	// olursa ekrana hata mesajı yazdırıyor ve hatanın detayınıda yazdırıyor ardından işlemi sonlandırıyor. Yani kısaca bu parça web sayfasının
	// içeriğini alıyor ve hata olursa güvenli şekilde duruyor.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("HTML okunurken hata olustu :", err)
		return
	}

	// bu kodun yaptığı şeyde şu: Web sayfasından okunnan içeriği body alıp "HtmlBilgisi.html" adlı bir dosyaya yazıyor. Dosya yazma sırasında bir hata olursa
	// ekrana hata mesajını ve hatanın detayını basıyor, ardından işlemi sonlandırıyor Yani kısaca, bu parça sitenin HTML bilgisini bilgisayara kaydediyor
	// ve hata olursa güvenli bir şekilde duruyor
	err = os.WriteFile("HtmlBilgisi.html", body, 0644)
	if err != nil {
		fmt.Println("HTML dosyaya yazilamadi :", err)
		return
	}

	fmt.Println("HTML icerigi HtmlBilgisi.html dosyasına kaydedildi") // kaydedildiğine dair bilgilendirme mesajı

	// Chromedp gerçek bir tarayıcıyı arka planda kontrol eder
	// Bu yüzden sayfanın birebir görüntüsü alınabilir

	/* Bu kod öncelikle Chrome tarayıcısı ile yeni bir çalışma bağlamı oluşturuyor. Bu bağlam, tarayıcıyla yapılacak otomasyon işlemlerinin (sayfa açma, element bulma, tıklama vb.)
	hangi ortamda çalışacağını belirliyor. Fonksiyon aynı zamanda cancel adında bir iptal fonksiyonu döndürüyor. defer cancel() ifadesi sayesinde, içinde bulunduğun fonksiyon
	bittiğinde bu iptal fonksiyonu otomatik çalışıyor ve açılan tarayıcı bağlamı kapatılıyor. Yani kısaca , bu
	kod Chrome otomasyonu için güvenli bir çalışma ortamı açıyor ve iş bitince kaynakları düzgün şekilde serbest bırakıyor. */
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	/* Programın çalışması sırasında belirli bir işlemin en fazla 30 saniye sürmesine izin verilir. Eğer bu süre aşılırsa, işlem otomatik olarak durdurulur.
	mevcut contexten yeni bir context üretilir ve bu context zaman kontrolünü üstlenir. Fonksiyon normal şekilde sona erse bile bu context in düzgün biçimde kapatılmasını sağlar
	Böylece hem işlemlerin gereksiz yere uzaması engellenir hem de sistem kaynakları boş yere harcanmaz. Kısacası bu yapı, uzun sürebilecek işlemleri kontrollü
	güvenli ve temiz bir şekilde yönetmek için kullanılır. */
	ctx, cancel = context.WithTimeout(ctx, time.Duration(*timeout)*time.Second)
	defer cancel()

	/* bu satırda screenshot adında bir değişken tanımlanır ve bu değişken byte türünde bir veri dizisini tutmak için kullanılır. []byte yapısı genellikle dosya içerikleri, ağdan gelen
	veriler veya bir ekran görüntüsü gibi ham verileri saklamak için tercih edilir. Bu değişken başlangıçta boş durumdadır ve daha sonra programın ilerleyen aşamalarında ekran görüntüsünün
	byte verileriyle doldurulması amaçlanır. Bu sayede alınan görüntü, dosyaya yazılabilir, ağ üzerinden gönderilebilir ya da başka işlemlerde kullanılabilir.
	*/
	var screenshot []byte

	/* chromedp kütüphanesi kullanılarak bir web sayfasına gidilmesini ve sayfanın ekran görüntüsünün alınmasını sağlar. verilen ctx context i içinde belirtilen adımları sırasıyla çalıştırır
	ilk olarak hedef web adresine gidilir ardından sayfanın tamamen yüklenmesi için kısa bir bekleme süresi tanımlanır. Son adımda ise sayfanın tamamının ekran görüntüsü alınır ve bu görüntü
	daha önce tanımlanan screenshot değişkenine byte formatında kaydedilir. tüm bu işlemler sırasında oluşabilecek hatalar err değişkeni üzerinden kontrol edilir.
	*/
	err = chromedp.Run(ctx,
		chromedp.Navigate(*url),
		chromedp.Sleep(2*time.Second), // Sayfanın tamamen yüklenmesi için kısa bekleme
		chromedp.FullScreenshot(&screenshot, 90),
	)

	/* ekran görüntüsü alma işlemi sırasında herhangi bir hata oluşup oluşmadığını kontrol etmek için kullanılır. err değişkeni boş değilse yani işlem sırasında bir sorun yaşanmamışsa, kullanıcıya durumu
	bildiren açıklayıcı bir hata mesajı ekrana yazdırılır. Ardından return ifadesi ile fonksiyonun çalışması durdurulur ve hatalı bir durumda programın devam etmesi engellenir. Bu yapı sayesinde program
	daha güvenili ve kontrollü çalışır ve oluşabilecek sorunlar sessizce geçilmek yerine açık bir şekilde raporlanmış olur.
	*/
	if err != nil {
		fmt.Println("Ekran goruntusu alinirken hata olustu :", err)
		return
	}

	/* daha önce alınan ekran görüntüsünü dosya olarak diske kaydetmek için kullanılır. screenshot değişkeninde tutulan byte verileri "EkranGoruntusu.png" adlı dosyaya yazılır. 0644 izin değeri,
	dosyanın sahibi tarafından okunup yazılabilmesini, diğer kullanıcılar tarafından ise sadece okunabilmesini sağlar. yazma işlemi sırasında bir hata oluşursa, bu durum kontrol edilir ce kullanıcıya
	ekran görüntüsünün dosyaya kaydedilmediğini belirten bir hata mesajı gösterilir. Ardından return ile fonksiyon sonlandırılarak hatalı bir durumda programın devam etmesi engellenir.
	*/
	err = os.WriteFile("EkranGoruntusu.png", screenshot, 0644)
	if err != nil {
		fmt.Println("Ekran goruntusu dosyaya yazilamadi :", err)
		return
	}

	fmt.Println("Ekran goruntusu screenshot.png olarak kaydedildi")
	fmt.Println("Program basariyla tamamlandi")
}
