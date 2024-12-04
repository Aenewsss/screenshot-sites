package main

import (
	"fmt"
	"log"
	// "strings"
	"time"

	"github.com/playwright-community/playwright-go"
)

func main() {
	// Inicializar o Playwright
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("Erro ao iniciar o Playwright: %v", err)
	}
	defer pw.Stop()

	// Iniciar o navegador no modo visível (não-headless)
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false), // Mantenha o navegador visível
	})
	if err != nil {
		log.Fatalf("Erro ao iniciar o navegador: %v", err)
	}
	defer browser.Close()

	// Criar um novo contexto de navegação
	context, err := browser.NewContext()
	if err != nil {
		log.Fatalf("Erro ao criar contexto: %v", err)
	}

	// Abrir uma nova página
	page, err := context.NewPage()
	if err != nil {
		log.Fatalf("Erro ao abrir nova página: %v", err)
	}

	// Navegar para o site desejado
	// url := "https://signmeupcustoms.com/"
	url := "https://aenamartinelli.com.br"
	// url := "https://www.apple.com/br/shop/buy-iphone"
	if _, err := page.Goto(url, playwright.PageGotoOptions{WaitUntil: playwright.WaitUntilStateNetworkidle}); err != nil {
		log.Fatalf("Erro ao navegar para o site: %v", err)
	}

	// Forçar carregamento de todas as imagens
	if err := loadAllImages(page); err != nil {
		log.Fatalf("Erro ao carregar imagens manualmente: %v", err)
	}

	// Aguardar explicitamente o carregamento das imagens
	// if err := waitForImagesToLoad(page); err != nil {
	// 	log.Fatalf("Erro ao aguardar carregamento de imagens: %v", err)
	// }

	// Realizar scroll até o final da página
	if err := scrollToBottom(page); err != nil {
		log.Fatalf("Erro ao realizar scroll: %v", err)
	}

	// Capturar a screenshot da página inteira
	_, err = page.Screenshot(playwright.PageScreenshotOptions{
		Path:     playwright.String("screenshot.png"),
		FullPage: playwright.Bool(true),
	})
	if err != nil {
		log.Fatalf("Erro ao tirar screenshot da página inteira: %v", err)
	}

	log.Println("Screenshot da página inteira capturada com sucesso.")

	// Aguarde o usuário para fechar o navegador
	fmt.Println("Navegador aberto. Pressione ENTER para encerrar...")
	fmt.Scanln()
}

// func removeHTTPPrefix(url string) string {
// 	url = strings.TrimPrefix(url, "http://")
// 	url = strings.TrimPrefix(url, "https://")
// 	return url
// }

// scrollToBottom faz scroll até o final da página
func scrollToBottom(page playwright.Page) error {
	const delayBetweenScrolls = 300 * time.Millisecond // Tempo entre eventos de scroll

	var previousDistance float64
	firstStep := true

	for {
		// Obter a altura atual do documento
		distanceFromTopInt, err := page.Evaluate("() => window.scrollY")
		if err != nil {
			return fmt.Errorf("erro ao obter altura total do documento: %w", err)
		}

		// Garantir que o valor seja um float64
		currentDistance, err := convertToFloat64(distanceFromTopInt)
		fmt.Printf("currentDistance %f\n", currentDistance)
		if err != nil {
			return fmt.Errorf("erro ao converter altura do documento: %w", err)
		}

		// Verificar se o scroll parou de aumentar
		if (currentDistance == previousDistance) && (previousDistance != 0 || !firstStep) {
			// Se a altura não mudou, chegamos ao final da página e inserimos um timeout antes do print
			time.Sleep(1 * time.Second)
			page.Evaluate(`
				const arrFiltered = Array.from(document.querySelectorAll('*')).filter(el => el.style.opacity != '1' && el.style.opacity != '');
				arrFiltered.forEach(el => {
					document.getElementsByClassName(el.className)[0].style.opacity = 1;
                	document.getElementsByClassName(el.className)[0].style.top = 'unset'
				})
			`)
			break
		}

		// Simular pressionamento da tecla "Page Down"
		if err := page.Keyboard().Press("PageDown"); err != nil {
			return fmt.Errorf("erro ao pressionar tecla PageDown: %w", err)
		}

		// Atualizar a altura anterior
		previousDistance = currentDistance

		firstStep = false

		// Aguardar para carregar novos conteúdos
		time.Sleep(delayBetweenScrolls)
	}

	return nil
}

// convertToFloat64 converte um valor interface{} para float64
func convertToFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	default:
		return 0, fmt.Errorf("tipo inesperado: %T", v)
	}
}

// loadAllImages força o carregamento manual de todas as imagens
func loadAllImages(page playwright.Page) error {
	script := `
		document.querySelectorAll("img").forEach(img => {
			if (!img.complete || img.naturalWidth === 0) {
				const src = img.getAttribute("data-src") || img.src;
				img.src = src;
			}
		});
	`
	_, err := page.Evaluate(script)
	return err
}

// waitForImagesToLoad aguarda o carregamento completo de todas as imagens
// func waitForImagesToLoad(page playwright.Page) error {
// 	script := `
// 		new Promise(resolve => {
// 			const images = Array.from(document.images);
// 			let loaded = 0;
// 			images.forEach(img => {
// 				if (img.complete && img.naturalWidth > 0) {
// 					loaded++;
// 				} else {
// 					img.onload = img.onerror = () => {
// 						loaded++;
// 						if (loaded === images.length) resolve();
// 					};
// 				}
// 			});
// 			if (loaded === images.length) resolve();
// 		});
// 	`
// 	_, err := page.Evaluate(script)
// 	return err
// }
