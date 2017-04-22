# Mystem Wrapper
Go wrapper of Yandex Mystem (only russian)

https://tech.yandex.ru/mystem/

# 
Данный пакет позволяет запускать Mystem из golang'a и получать выходной текст, где каждое слово приведено в свою исходную форму.
Это нужно для подготовки датасетов для различных семантических анализов.

Например, предыдущие предложения после обработки будут выглядеть так:
```
данный пакет позволять запускать Mystem из golang'a и получать выходной текст, 
где каждый слово приводить в свой исходный форма.
это нужно для подготовка датасет для различный семантический анализ.
```

Пример использования:
```
package main

import (
	"github.com/un000/mystem-wrapper"
	"fmt"
)

func main() {
  // create wrapper
	m := mystem.New("../../bin/mystem", []string{"-d"}) // -d is for context homonymization

	sourceTexts := []string{`Данный пакет позволяет запускать Mystem из golang'a и получать выходной текст,
    где каждое слово приведено в свою исходную форму.
		Это нужно для подготовки датасетов для различных семантических анализов.`}

	transformedTexts, err := m.Transform(sourceTexts)
	if err != nil {
		fmt.Print(err)
	}

	for i := range transformedTexts {
		fmt.Println(transformedTexts[i])
	}
}
  ```
