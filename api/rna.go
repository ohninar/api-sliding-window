package api

import (
	"fmt"
	"math"

	"github.com/mitsuse/matrix-go/dense"
)

var (
	//theta1 = dense.New(512, 901)(Theta1...)
	//theta2 = dense.New(62, 513)(Theta2...)
	limiar = 0.5
)

/*
[1] - X = matriz(1x900) //1x900
[2] - X1 = matriz(1x901) // [X 1] vetorCaracteristica adicionado o BIAS; 1x901
[3] - X11 = X1 * Theta1'  //1x901 * (512x901)'  => 1x901 * 901x512 =>resultdo = 1x512
[4] - X11sigmoid = sigmoid(X11); //1x512
%SIGMOID Compute sigmoid functoon
%   J = SIGMOID(z) computes the sigmoid of z.
% g = 1.0 ./ (1.0 + exp(-z));

[5] - X2 = [X11sigmoid 1]; //1x513
[6] - X22 = X2 * Theta2'; // 1x513 * (62x513)' => 1x513 * 513x62 => 1x62
[7] -  = sigmoid(X22); // 1x62
%normalizar a saida
8 - MatrizPred = X22sigmoid(i)/sum(X22sigmoid); // 1x62
[9] - [LINHA, COLUNA] = max(MatrizPred)
LINHA = 0.76667
COLUNA = 57
[10] - aplica a limiar LINHA > limiar? sim, LINHA = LINHA;
11 - Conclusão a nova amostra foi classificada como sendo da
classe 57 com "precisão" de 76%

faz de conta um novo exemplo
// ABA..ATE
9 - [LINHA, COLUNA] = max(MatrizPred)
LINHA = 0.36667
COLUNA = 15
10 - aplica a limiar LINHA > limiar? não, LINHA = 0;
11 - Conclusão a nova amostra foi classificada como sendo da
classe 15 com "precisão" de 36%
*/

func sorterImage(matrix []float64) string {
	newMatrix := setMatrix(matrix, 1, 901)
	newMatrix = setBIAS(newMatrix)
	newMatrix = multiplicationMatrix(newMatrix, transposeMatrix(theta1))
	matrix = sigmoid(newMatrix)

	newMatrix = setMatrix(matrix, 1, 513)
	newMatrix = setBIAS(newMatrix)
	newMatrix = multiplicationMatrix(newMatrix, transposeMatrix(theta2))
	matrix = sigmoid(newMatrix)
	newMatrix = setMatrix(matrix, 1, 62)
	fmt.Println(newMatrix)
	element, _, column := newMatrix.Max()
	fmt.Println(element, column)
	if checkLimiar(element, limiar) {
		return represatation(column)
	}
	return ""
}

func represatation(column int) string {
	switch {
	case 1 == column:
		return "0"
	case 2 == column:
		return "1"
	case 3 == column:
		return "2"
	case 4 == column:
		return "3"
	case 5 == column:
		return "4"
	case 6 == column:
		return "5"
	case 7 == column:
		return "6"
	case 8 == column:
		return "7"
	case 9 == column:
		return "8"
	case 10 == column:
		return "9"
	case 11 == column:
		return "A"
	case 12 == column:
		return "B"
	case 13 == column:
		return "C"
	case 14 == column:
		return "D"
	case 15 == column:
		return "E"
	case 16 == column:
		return "F"
	case 17 == column:
		return "G"
	case 18 == column:
		return "H"
	case 19 == column:
		return "I"
	case 20 == column:
		return "J"
	case 21 == column:
		return "K"
	case 22 == column:
		return "L"
	case 23 == column:
		return "M"
	case 24 == column:
		return "N"
	case 25 == column:
		return "O"
	case 26 == column:
		return "P"
	case 27 == column:
		return "Q"
	case 28 == column:
		return "R"
	case 29 == column:
		return "S"
	case 30 == column:
		return "T"
	case 31 == column:
		return "U"
	case 32 == column:
		return "V"
	case 33 == column:
		return "W"
	case 34 == column:
		return "X"
	case 35 == column:
		return "Y"
	case 36 == column:
		return "Z"
	case 37 == column:
		return "a"
	case 38 == column:
		return "b"
	case 39 == column:
		return "c"
	case 40 == column:
		return "d"
	case 41 == column:
		return "e"
	case 42 == column:
		return "f"
	case 43 == column:
		return "g"
	case 44 == column:
		return "h"
	case 45 == column:
		return "i"
	case 46 == column:
		return "j"
	case 47 == column:
		return "k"
	case 48 == column:
		return "l"
	case 49 == column:
		return "m"
	case 50 == column:
		return "n"
	case 51 == column:
		return "o"
	case 52 == column:
		return "p"
	case 53 == column:
		return "q"
	case 54 == column:
		return "r"
	case 55 == column:
		return "s"
	case 56 == column:
		return "t"
	case 57 == column:
		return "u"
	case 58 == column:
		return "v"
	case 59 == column:
		return "w"
	case 60 == column:
		return "x"
	case 61 == column:
		return "y"
	case 62 == column:
		return "z"
	}
	return ""
}

//fazer um for da matri float e inserir no dense matrix e não o aucontrario
func setMatrix(matrix []float64, row int, column int) *dense.Matrix {
	newMatrix := dense.Zeros(row, column)
	for i := 0; i < len(matrix); i++ {
		newMatrix.Update(row-1, i, matrix[i])
	}
	return newMatrix
}

func checkLimiar(element float64, limiar float64) bool {
	if element > limiar {
		return true
	}
	return false
}

func sigmoid(matrix *dense.Matrix) []float64 {
	cursor := matrix.All()
	newMatrix := []float64{}

	for cursor.HasNext() {
		element, _, _ := cursor.Get()
		newElement := 1 / (1 + math.Exp(-element))
		newMatrix = append(newMatrix, newElement)
	}

	return newMatrix
}

func multiplicationMatrix(matrix *dense.Matrix, theta *dense.Matrix) *dense.Matrix {
	newMatrix := matrix.Multiply(theta)
	return dense.Convert(newMatrix)
}

func setBIAS(matrix *dense.Matrix) *dense.Matrix {
	row := matrix.Rows()
	column := matrix.Columns()
	matrix.Update(row-1, column-1, 1.0)
	return matrix
}

func transposeMatrix(matrix *dense.Matrix) *dense.Matrix {
	newMatrix := dense.Zeros(matrix.Columns(), matrix.Rows())
	cursor := matrix.All()

	for cursor.HasNext() {
		element, row, column := cursor.Get()
		newMatrix.Update(column, row, element)
	}

	return newMatrix
}
