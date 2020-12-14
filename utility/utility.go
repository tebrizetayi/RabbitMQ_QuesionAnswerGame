package utility

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"reflect"
	"strings"
)

const (
	//AnswersQ is queue for answers(incoming)
	AnswersQ = "answers"
	//QuestionQ is queue for questions(outgoing)
	QuestionQ = "questions"
)

//CreateExam creates question together with answer
func CreateExam() map[string]string {
	fileName := "../QuestionsAnswers.txt"
	file, err := os.Open(fileName)
	FailOnError(err, "Could not open the file")
	defer file.Close()

	questions := make(map[string]string)
	scanner := bufio.NewScanner(file)
	question := ""
	answer := ""
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Index(line, "Answer") == 0 {
			answer = line
			questions[question] = answer
			question = ""
			answer = ""

		} else {
			question += line
		}

	}

	return questions
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func PickRandomKey(dict map[string]string) string {

	keys := reflect.ValueOf(dict).MapKeys()
	return fmt.Sprintf("%v", keys[rand.Intn(len(keys))])

}
