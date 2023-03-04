package gpt

func formatAnswer(answer string) string {
	for len(answer) > 0 {
		if answer[:1] == "\n" || answer[0] == ' ' {
			answer = answer[1:]
		} else {
			break
		}
	}
	return answer
}

func GetLen(s []string) int {
	count := 0
	for _, v := range s {
		count += len(v)
	}
	return count
}

func GetMaxSubset(s []string, max int) []string {
	for i := 0; i <= len(s); i++ {
		if GetLen(s[i:]) <= max {
			return s[i:]
		}
	}
	return []string{}
}
