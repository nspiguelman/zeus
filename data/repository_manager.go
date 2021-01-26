package data

type RepositoryManager struct {
	KahootRepository   *KahootRepository
	UserRepository     *UserRepository
	QuestionRepository *QuestionRepository
	AnswerRepository   *AnswerRepository
}

func NewRepositoryManager(dbData *Data) *RepositoryManager {
	return &RepositoryManager{
		KahootRepository: &KahootRepository{
			Data: data,
		},
		UserRepository: &UserRepository{
			Data: dbData,
		},
		QuestionRepository: &QuestionRepository{
			Data: dbData,
		},
		AnswerRepository: &AnswerRepository{
			Data: dbData,
		},
	}
}
