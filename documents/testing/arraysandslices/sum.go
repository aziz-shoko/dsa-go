package arraysandslices

func Sum(a []int) int {
	sum := 0
	for _, i := range a {
		sum += i	
	}
	return sum
}

func SumAll(numbersToSum ...[]int) []int {
	var sums []int
	for _, numbers := range numbersToSum {
		sums = append(sums, Sum(numbers))
	}

	return sums
}

func SumAllTails(tailsToSum ...[]int) []int {
	var sums []int
	for _, numbers := range tailsToSum {
		if len(numbers) > 0 {
			tailSum := numbers[1:]
			sums = append(sums, Sum(tailSum))
		} else {
			sums = append(sums, 0)
		}
	}
	return sums
}