package test

import (
	"bufio"
	"goUmbral/utils"
	"os"
	"strconv"
	"strings"
	"testing"
)

func BenchmarkGzipCompressRandom(b *testing.B) {
	expected, _ := os.ReadFile("random1.txt")
	utils.GzipCompress(expected)
}

func BenchmarkGzipCompressRegular(b *testing.B) {
	expected, _ := os.ReadFile("regular1.txt")
	utils.GzipCompress(expected)
}

func BenchmarkGzipCompressSame(b *testing.B) {
	expected, _ := os.ReadFile("same1.txt")
	utils.GzipCompress(expected)
}

func BenchmarkGzipCompressCA(b *testing.B) {
	expected, _ := os.ReadFile("CA_fameleBirth.txt")
	utils.GzipCompress(expected)
}

func BenchmarkGzipCompressSun(b *testing.B) {
	expected, _ := os.ReadFile("sunspot.txt")
	utils.GzipCompress(expected)
}

func BenchmarkGzipUncompressRandom(b *testing.B) {
	b.StopTimer()
	expected, _ := os.ReadFile("random1.txt")
	data1 := utils.GzipCompress(expected)
	b.StartTimer()
	utils.GzipUnCompress(data1)
}

func BenchmarkGzipUncompressRegular(b *testing.B) {
	b.StopTimer()
	expected, _ := os.ReadFile("regular1.txt")
	data1 := utils.GzipCompress(expected)
	b.StartTimer()
	utils.GzipUnCompress(data1)
}

func BenchmarkGzipUncompressSame(b *testing.B) {
	b.StopTimer()
	expected, _ := os.ReadFile("same1.txt")
	data1 := utils.GzipCompress(expected)
	b.StartTimer()
	utils.GzipUnCompress(data1)
}

func BenchmarkGzipUncompressCA(b *testing.B) {
	b.StopTimer()
	expected, _ := os.ReadFile("CA_fameleBirth.txt")
	data1 := utils.GzipCompress(expected)
	b.StartTimer()
	utils.GzipUnCompress(data1)
}

func BenchmarkGzipUncompressSun(b *testing.B) {
	b.StopTimer()
	expected, _ := os.ReadFile("sunspot.txt")
	data1 := utils.GzipCompress(expected)
	b.StartTimer()
	utils.GzipUnCompress(data1)
}

func BenchmarkGorillaCompressRandom(b *testing.B) {
	type Point struct {
		V float64
		T uint32
	}
	expected := []Point{}
	fRead, _ := os.Open("random1.txt")
	defer fRead.Close()
	fs := bufio.NewScanner(fRead)
	for fs.Scan() {
		strline := fs.Text()
		strlineSplict := strings.Split(strline, " ")
		value, _ := strconv.ParseFloat(strlineSplict[0], 64)
		time1, _ := strconv.ParseUint(strlineSplict[1], 10, 32)
		time := uint32(time1)
		expected = append(expected, Point{value, time})
	}
	s1 := utils.New(expected[0].T)
	for _, p := range expected {
		s1.Push(p.T, p.V)
	}
	s1.MarshalBinary()
}

func BenchmarkGorillaCompressRegular(b *testing.B) {
	type Point struct {
		V float64
		T uint32
	}
	expected := []Point{}
	fRead, _ := os.Open("regular1.txt")
	defer fRead.Close()
	fs := bufio.NewScanner(fRead)
	for fs.Scan() {
		strline := fs.Text()
		strlineSplict := strings.Split(strline, " ")
		value, _ := strconv.ParseFloat(strlineSplict[0], 64)
		time1, _ := strconv.ParseUint(strlineSplict[1], 10, 32)
		time := uint32(time1)
		expected = append(expected, Point{value, time})
	}
	s1 := utils.New(expected[0].T)
	for _, p := range expected {
		s1.Push(p.T, p.V)
	}
	s1.MarshalBinary()
}

func BenchmarkGorillaCompressSame(b *testing.B) {
	type Point struct {
		V float64
		T uint32
	}
	expected := []Point{}
	fRead, _ := os.Open("same1.txt")
	defer fRead.Close()
	fs := bufio.NewScanner(fRead)
	for fs.Scan() {
		strline := fs.Text()
		strlineSplict := strings.Split(strline, " ")
		value, _ := strconv.ParseFloat(strlineSplict[0], 64)
		time1, _ := strconv.ParseUint(strlineSplict[1], 10, 32)
		time := uint32(time1)
		expected = append(expected, Point{value, time})
	}
	s1 := utils.New(expected[0].T)
	for _, p := range expected {
		s1.Push(p.T, p.V)
	}
	s1.MarshalBinary()
}

func BenchmarkGorillaCompressCA(b *testing.B) {
	type Point struct {
		V float64
		T uint32
	}
	expected := []Point{}
	fRead, _ := os.Open("CA_fameleBirth.txt")
	defer fRead.Close()
	fs := bufio.NewScanner(fRead)
	for fs.Scan() {
		strline := fs.Text()
		strlineSplict := strings.Split(strline, " ")
		value, _ := strconv.ParseFloat(strlineSplict[0], 64)
		time1, _ := strconv.ParseUint(strlineSplict[1], 10, 32)
		time := uint32(time1)
		expected = append(expected, Point{value, time})
	}
	s1 := utils.New(expected[0].T)
	for _, p := range expected {
		s1.Push(p.T, p.V)
	}
	s1.MarshalBinary()
}

func BenchmarkGorillaCompressSun(b *testing.B) {
	type Point struct {
		V float64
		T uint32
	}
	expected := []Point{}
	fRead, _ := os.Open("sunspot.txt")
	defer fRead.Close()
	fs := bufio.NewScanner(fRead)
	for fs.Scan() {
		strline := fs.Text()
		strlineSplict := strings.Split(strline, " ")
		value, _ := strconv.ParseFloat(strlineSplict[0], 64)
		time1, _ := strconv.ParseUint(strlineSplict[1], 10, 32)
		time := uint32(time1)
		expected = append(expected, Point{value, time})
	}
	s1 := utils.New(expected[0].T)
	for _, p := range expected {
		s1.Push(p.T, p.V)
	}
	s1.MarshalBinary()
}

func BenchmarkGorillaUncompressRandom(b *testing.B) {
	b.StopTimer()
	type Point struct {
		V float64
		T uint32
	}
	expected := []Point{}
	fRead, _ := os.Open("random1.txt")
	defer fRead.Close()
	fs := bufio.NewScanner(fRead)
	for fs.Scan() {
		strline := fs.Text()
		strlineSplict := strings.Split(strline, " ")
		value, _ := strconv.ParseFloat(strlineSplict[0], 64)
		time1, _ := strconv.ParseUint(strlineSplict[1], 10, 32)
		time := uint32(time1)
		expected = append(expected, Point{value, time})
	}
	s1 := utils.New(expected[0].T)
	for _, p := range expected {
		s1.Push(p.T, p.V)
	}
	data1, _ := s1.MarshalBinary()
	os.WriteFile("compress.txt", data1, 0644)
	b.StartTimer()
	data2, _ := os.ReadFile("compress.txt")
	s2 := utils.New(uint32(1))
	s2.UnmarshalBinary(data2)
	os.Remove("compress.txt")
}

func BenchmarkGorillaUncompressRegular(b *testing.B) {
	b.StopTimer()
	type Point struct {
		V float64
		T uint32
	}
	expected := []Point{}
	fRead, _ := os.Open("regular1.txt")
	defer fRead.Close()
	fs := bufio.NewScanner(fRead)
	for fs.Scan() {
		strline := fs.Text()
		strlineSplict := strings.Split(strline, " ")
		value, _ := strconv.ParseFloat(strlineSplict[0], 64)
		time1, _ := strconv.ParseUint(strlineSplict[1], 10, 32)
		time := uint32(time1)
		expected = append(expected, Point{value, time})
	}
	s1 := utils.New(expected[0].T)
	for _, p := range expected {
		s1.Push(p.T, p.V)
	}
	data1, _ := s1.MarshalBinary()
	os.WriteFile("compress.txt", data1, 0644)
	b.StartTimer()
	data2, _ := os.ReadFile("compress.txt")
	s2 := utils.New(uint32(1))
	s2.UnmarshalBinary(data2)
	os.Remove("compress.txt")
}

func BenchmarkGorillaUncompressSame(b *testing.B) {
	b.StopTimer()
	type Point struct {
		V float64
		T uint32
	}
	expected := []Point{}
	fRead, _ := os.Open("same1.txt")
	defer fRead.Close()
	fs := bufio.NewScanner(fRead)
	for fs.Scan() {
		strline := fs.Text()
		strlineSplict := strings.Split(strline, " ")
		value, _ := strconv.ParseFloat(strlineSplict[0], 64)
		time1, _ := strconv.ParseUint(strlineSplict[1], 10, 32)
		time := uint32(time1)
		expected = append(expected, Point{value, time})
	}
	s1 := utils.New(expected[0].T)
	for _, p := range expected {
		s1.Push(p.T, p.V)
	}
	data1, _ := s1.MarshalBinary()
	os.WriteFile("compress.txt", data1, 0644)
	b.StartTimer()
	data2, _ := os.ReadFile("compress.txt")
	s2 := utils.New(uint32(1))
	s2.UnmarshalBinary(data2)
	os.Remove("compress.txt")
}

func BenchmarkGorillaUncompressCA(b *testing.B) {
	b.StopTimer()
	type Point struct {
		V float64
		T uint32
	}
	expected := []Point{}
	fRead, _ := os.Open("CA_fameleBirth.txt")
	defer fRead.Close()
	fs := bufio.NewScanner(fRead)
	for fs.Scan() {
		strline := fs.Text()
		strlineSplict := strings.Split(strline, " ")
		value, _ := strconv.ParseFloat(strlineSplict[0], 64)
		time1, _ := strconv.ParseUint(strlineSplict[1], 10, 32)
		time := uint32(time1)
		expected = append(expected, Point{value, time})
	}
	s1 := utils.New(expected[0].T)
	for _, p := range expected {
		s1.Push(p.T, p.V)
	}
	data1, _ := s1.MarshalBinary()
	os.WriteFile("compress.txt", data1, 0644)
	b.StartTimer()
	data2, _ := os.ReadFile("compress.txt")
	s2 := utils.New(uint32(1))
	s2.UnmarshalBinary(data2)
	os.Remove("compress.txt")
}

func BenchmarkGorillaUncompressSun(b *testing.B) {
	b.StopTimer()
	type Point struct {
		V float64
		T uint32
	}
	expected := []Point{}
	fRead, _ := os.Open("sunspot.txt")
	defer fRead.Close()
	fs := bufio.NewScanner(fRead)
	for fs.Scan() {
		strline := fs.Text()
		strlineSplict := strings.Split(strline, " ")
		value, _ := strconv.ParseFloat(strlineSplict[0], 64)
		time1, _ := strconv.ParseUint(strlineSplict[1], 10, 32)
		time := uint32(time1)
		expected = append(expected, Point{value, time})
	}
	s1 := utils.New(expected[0].T)
	for _, p := range expected {
		s1.Push(p.T, p.V)
	}
	data1, _ := s1.MarshalBinary()
	os.WriteFile("compress.txt", data1, 0644)
	b.StartTimer()
	data2, _ := os.ReadFile("compress.txt")
	s2 := utils.New(uint32(1))
	s2.UnmarshalBinary(data2)
	os.Remove("compress.txt")
}

func TestGzipCompressRateRandom(t *testing.T) {
	expected, _ := os.ReadFile("random2.txt")
	t.Log("文件原始大小:", len(expected))
	data := utils.GzipCompress(expected)
	t.Log("文件压缩大小:", len(data))
	rate := (float64)(len(data)) / float64(len(expected)) * 100
	t.Log("压缩率:", rate, "%")
}

func TestGzipCompressRateRegular(t *testing.T) {
	expected, _ := os.ReadFile("regular2.txt")
	t.Log("文件原始大小:", len(expected))
	data := utils.GzipCompress(expected)
	t.Log("文件压缩大小:", len(data))
	rate := (float64)(len(data)) / float64(len(expected)) * 100
	t.Log("压缩率:", rate, "%")
}

func TestGzipCompressRateSame(t *testing.T) {
	expected, _ := os.ReadFile("same2.txt")
	t.Log("文件原始大小:", len(expected))
	data := utils.GzipCompress(expected)
	t.Log("文件压缩大小:", len(data))
	rate := (float64)(len(data)) / float64(len(expected)) * 100
	t.Log("压缩率:", rate, "%")
}

func TestGzipCompressRateCA(t *testing.T) {
	expected, _ := os.ReadFile("CA_fameleBirth.txt")
	t.Log("文件原始大小:", len(expected))
	data := utils.GzipCompress(expected)
	t.Log("文件压缩大小:", len(data))
	rate := (float64)(len(data)) / float64(len(expected)) * 100
	t.Log("压缩率:", rate, "%")
}

func TestGzipCompressRateSun(t *testing.T) {
	expected, _ := os.ReadFile("sunspot.txt")
	t.Log("文件原始大小:", len(expected))
	data := utils.GzipCompress(expected)
	t.Log("文件压缩大小:", len(data))
	rate := (float64)(len(data)) / float64(len(expected)) * 100
	t.Log("压缩率:", rate, "%")
}

func TestGorillaCompressRateRandom(t *testing.T) {
	type Point struct {
		V float64
		T uint32
	}
	expected := []Point{}
	fRead, _ := os.Open("random2.txt")
	defer fRead.Close()
	fs := bufio.NewScanner(fRead)
	for fs.Scan() {
		strline := fs.Text()
		strlineSplict := strings.Split(strline, " ")
		value, _ := strconv.ParseFloat(strlineSplict[0], 64)
		time1, _ := strconv.ParseUint(strlineSplict[1], 10, 32)
		time := uint32(time1)
		expected = append(expected, Point{value, time})
	}
	s1 := utils.New(expected[0].T)
	for _, p := range expected {
		s1.Push(p.T, p.V)
	}
	data, _ := s1.MarshalBinary()

	origin, _ := os.ReadFile("random2.txt")
	t.Log("文件原始大小:", len(origin))
	t.Log("文件压缩大小:", len(data))
	rate := (float64)(len(data)) / float64(len(origin)) * 100
	t.Log("压缩率:", rate, "%")
}

func TestGorillaCompressRateRegular(t *testing.T) {
	type Point struct {
		V float64
		T uint32
	}
	expected := []Point{}
	fRead, _ := os.Open("regular2.txt")
	defer fRead.Close()
	fs := bufio.NewScanner(fRead)
	for fs.Scan() {
		strline := fs.Text()
		strlineSplict := strings.Split(strline, " ")
		value, _ := strconv.ParseFloat(strlineSplict[0], 64)
		time1, _ := strconv.ParseUint(strlineSplict[1], 10, 32)
		time := uint32(time1)
		expected = append(expected, Point{value, time})
	}
	s1 := utils.New(expected[0].T)
	for _, p := range expected {
		s1.Push(p.T, p.V)
	}
	data, _ := s1.MarshalBinary()

	origin, _ := os.ReadFile("regular2.txt")
	t.Log("文件原始大小:", len(origin))
	t.Log("文件压缩大小:", len(data))
	rate := (float64)(len(data)) / float64(len(origin)) * 100
	t.Log("压缩率:", rate, "%")
}

func TestGorillaCompressRateSame(t *testing.T) {
	type Point struct {
		V float64
		T uint32
	}
	expected := []Point{}
	fRead, _ := os.Open("same2.txt")
	defer fRead.Close()
	fs := bufio.NewScanner(fRead)
	for fs.Scan() {
		strline := fs.Text()
		strlineSplict := strings.Split(strline, " ")
		value, _ := strconv.ParseFloat(strlineSplict[0], 64)
		time1, _ := strconv.ParseUint(strlineSplict[1], 10, 32)
		time := uint32(time1)
		expected = append(expected, Point{value, time})
	}
	s1 := utils.New(expected[0].T)
	for _, p := range expected {
		s1.Push(p.T, p.V)
	}
	data, _ := s1.MarshalBinary()

	origin, _ := os.ReadFile("same2.txt")
	t.Log("文件原始大小:", len(origin))
	t.Log("文件压缩大小:", len(data))
	rate := (float64)(len(data)) / float64(len(origin)) * 100
	t.Log("压缩率:", rate, "%")
}

func TestGorillaCompressRateCA(t *testing.T) {
	type Point struct {
		V float64
		T uint32
	}
	expected := []Point{}
	fRead, _ := os.Open("CA_fameleBirth.txt")
	defer fRead.Close()
	fs := bufio.NewScanner(fRead)
	for fs.Scan() {
		strline := fs.Text()
		strlineSplict := strings.Split(strline, " ")
		value, _ := strconv.ParseFloat(strlineSplict[0], 64)
		time1, _ := strconv.ParseUint(strlineSplict[1], 10, 32)
		time := uint32(time1)
		expected = append(expected, Point{value, time})
	}
	s1 := utils.New(expected[0].T)
	for _, p := range expected {
		s1.Push(p.T, p.V)
	}
	data, _ := s1.MarshalBinary()

	origin, _ := os.ReadFile("CA_fameleBirth.txt")
	t.Log("文件原始大小:", len(origin))
	t.Log("文件压缩大小:", len(data))
	rate := (float64)(len(data)) / float64(len(origin)) * 100
	t.Log("压缩率:", rate, "%")
}

func TestGorillaCompressRateSun(t *testing.T) {
	type Point struct {
		V float64
		T uint32
	}
	expected := []Point{}
	fRead, _ := os.Open("sunspot.txt")
	defer fRead.Close()
	fs := bufio.NewScanner(fRead)
	for fs.Scan() {
		strline := fs.Text()
		strlineSplict := strings.Split(strline, " ")
		value, _ := strconv.ParseFloat(strlineSplict[0], 64)
		time1, _ := strconv.ParseUint(strlineSplict[1], 10, 32)
		time := uint32(time1)
		expected = append(expected, Point{value, time})
	}
	s1 := utils.New(expected[0].T)
	for _, p := range expected {
		s1.Push(p.T, p.V)
	}
	data, _ := s1.MarshalBinary()

	origin, _ := os.ReadFile("sunspot.txt")
	t.Log("文件原始大小:", len(origin))
	t.Log("文件压缩大小:", len(data))
	rate := (float64)(len(data)) / float64(len(origin)) * 100
	t.Log("压缩率:", rate, "%")
}
