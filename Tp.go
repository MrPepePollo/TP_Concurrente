package main

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

func RegresionLinealSecuencial(x []float64, y []float64) (float64, float64) {
	n := float64(len(x))
	var sumX, sumY, sumXY, sumX2 float64

	for i := 0; i < int(n); i++ {
		sumX += x[i]
		sumY += y[i]
		sumXY += x[i] * y[i]
		sumX2 += x[i] * x[i]
	}

	m := (n*sumXY - sumX*sumY) / (n*sumX2 - sumX*sumX)
	b := (sumY - m*sumX) / n

	return m, b
}

type Resultado struct {
	SumX  float64
	SumY  float64
	SumXY float64
	SumX2 float64
	mux   sync.Mutex
}

func ParteRegresion(x []float64, y []float64, wg *sync.WaitGroup, resultado *Resultado) {
	defer wg.Done()

	var sumX, sumY, sumXY, sumX2 float64

	for i := 0; i < len(x); i++ {
		sumX += x[i]
		sumY += y[i]
		sumXY += x[i] * y[i]
		sumX2 += x[i] * x[i]
	}

	resultado.mux.Lock()
	resultado.SumX += sumX
	resultado.SumY += sumY
	resultado.SumXY += sumXY
	resultado.SumX2 += sumX2
	resultado.mux.Unlock()
}

func RegresionLinealConcurrente(x []float64, y []float64, partes int) (float64, float64) {
	n := float64(len(x))
	partSize := len(x) / partes
	resultado := &Resultado{}

	var wg sync.WaitGroup

	for i := 0; i < partes; i++ {
		inicio := i * partSize
		fin := inicio + partSize
		if i == partes-1 {
			fin = len(x)
		}
		wg.Add(1)
		go ParteRegresion(x[inicio:fin], y[inicio:fin], &wg, resultado)
	}

	wg.Wait()

	m := (n*resultado.SumXY - resultado.SumX*resultado.SumY) / (n*resultado.SumX2 - resultado.SumX*resultado.SumX)
	b := (resultado.SumY - m*resultado.SumX) / n

	return m, b
}

func GenerarDatos(n int) ([]float64, []float64) {
	x := make([]float64, n)
	y := make([]float64, n)

	for i := 0; i < n; i++ {
		x[i] = float64(i)
		y[i] = float64(2*i + 3)
	}

	return x, y
}

func main() {
	n := 1000000 // Tamaño del conjunto de datos
	partes := 4  // Numero de partes para concurrencia

	x, y := GenerarDatos(n)

	const numPruebas = 1000
	numRep := 5 // Número de repeticiones

	for rep := 1; rep <= numRep; rep++ {
		// Calcular regresion lineal secuencial
		mSec, bSec := RegresionLinealSecuencial(x, y)
		fmt.Printf("Prueba %d - Regresión lineal secuencial: Pendiente=%.6f, Intercepto=%.6f\n", rep, mSec, bSec)

		// Calcular regresion lineal concurrente
		mCon, bCon := RegresionLinealConcurrente(x, y, partes)
		fmt.Printf("Prueba %d - Regresión lineal concurrente: Pendiente=%.6f, Intercepto=%.6f\n", rep, mCon, bCon)

		// Pruebas de tiempo para la implementacion secuencial
		tiemposSecuencial := make([]float64, numPruebas)

		for i := 0; i < numPruebas; i++ {
			inicio := time.Now()
			RegresionLinealSecuencial(x, y)
			fin := time.Since(inicio).Seconds()
			tiemposSecuencial[i] = fin
		}

		// Pruebas de tiempo para la implementacion concurrente
		tiemposConcurrente := make([]float64, numPruebas)

		for i := 0; i < numPruebas; i++ {
			inicio := time.Now()
			RegresionLinealConcurrente(x, y, partes)
			fin := time.Since(inicio).Seconds()
			tiemposConcurrente[i] = fin
		}

		// Calcular la media recortada para la implementación secuencial
		sort.Float64s(tiemposSecuencial)
		tiemposRecortadosSec := tiemposSecuencial[50 : len(tiemposSecuencial)-50]
		sumTiemposSec := 0.0
		for _, t := range tiemposRecortadosSec {
			sumTiemposSec += t
		}
		mediaRecortadaSec := sumTiemposSec / float64(len(tiemposRecortadosSec))

		// Calcular la media recortada para la implementación concurrente
		sort.Float64s(tiemposConcurrente)
		tiemposRecortadosCon := tiemposConcurrente[50 : len(tiemposConcurrente)-50]
		sumTiemposCon := 0.0
		for _, t := range tiemposRecortadosCon {
			sumTiemposCon += t
		}
		mediaRecortadaCon := sumTiemposCon / float64(len(tiemposRecortadosCon))

		fmt.Printf("Prueba %d - Media recortada del tiempo de ejecución (secuencial): %.6f segundos\n", rep, mediaRecortadaSec)
		fmt.Printf("Prueba %d - Media recortada del tiempo de ejecución (concurrente): %.6f segundos\n", rep, mediaRecortadaCon)
	}
}
