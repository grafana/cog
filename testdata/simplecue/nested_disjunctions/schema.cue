HelloMsg: {
    type: "hello"
    salutation: string
}

ByeMsg: {
    type: "bye"
    reason: string
}

QuestionMsg: {
    type: "question"
    question: string
}

AnswerMsg: {
    type: "answer"
    content: string
}

ChitChat: HelloMsg | ByeMsg
Message: ChitChat | QuestionMsg | AnswerMsg
