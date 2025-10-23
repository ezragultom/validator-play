package main

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"ezragultom/validator-play/apperror"
	"ezragultom/validator-play/constants"

	"github.com/go-playground/validator/v10"
)

type ProgramCreateRequest struct {
	Title              string                   `json:"title" validate:"required,min=4,max=255"`
	Description        string                   `json:"description" validate:"required"`
	CoverImage         string                   `json:"cover_image" validate:"required"`
	SupportingDocument string                   `json:"supporting_document"` //optional
	SubmissionDeadline string                   `json:"submission_deadline" validate:"required,date"`
	ProgramPeriod      string                   `json:"program_period" validate:"required"`
	ImplementationCost int64                    `json:"implementation_cost" validate:"gt=0"`
	Question           []ProgramQuestionRequest `json:"question" validate:"required,dive"`
}

type ProgramQuestionRequest struct {
	QuestionType string              `json:"question_type" validate:"required,oneof=short_text long_text date_time file_upload url_link option"`
	Label        string              `json:"label" validate:"required"`
	HelperText   string              `json:"helper_text"` //optional
	Attributes   ProgramQuestionAttr `json:"attributes"`
	ShowingOrder int                 `json:"showing_order" validate:"required,gte=0"`
	IsRequired   bool                `json:"is_required"`
}

type ProgramQuestionAttr struct {
	Placeholder     string   `json:"placeholder,omitempty" validate:"required_by_question_type_text"` //needed for question_type -> short_text long_text
	MaxChar         int      `json:"max_char,omitempty" validate:"gte=0"`                             //optional
	InputValidation string   `json:"input_validation,omitempty" validate:"required_if_short_text"`
	Format          string   `json:"format,omitempty" validate:"required_by_question_type_date_time"` //needed for question_type -> date_time
	OptionType      string   `json:"option_type,omitempty" validate:"required_by_question_type_option"`
	Choices         []string `json:"choices,omitempty" validate:"required_by_question_type_option_choices"`
}

func ValidateDate(fl validator.FieldLevel) bool {
	_, err := time.Parse(time.DateOnly, fl.Field().String())
	return err == nil
}

func ValidateQuestionAttributesTextType(fl validator.FieldLevel) bool {
	questionType := getQuestionType(fl)
	if questionType == constants.QuestionTypeShortText || questionType == constants.QuestionTypeLongText {
		return fl.Field().String() != ""
	}

	return true
}

func ValidateQuestionAttributesOptionType(fl validator.FieldLevel) bool {
	questionType := getQuestionType(fl)
	if questionType == constants.QuestionTypeOption {
		fmt.Println(fl.Field().String())
		fmt.Println(fl.Field().String() == constants.QuestionAttributeOptionTypeRadioButton || fl.Field().String() == constants.QuestionAttributeOptionTypeCheckbox)
		return fl.Field().String() == constants.QuestionAttributeOptionTypeRadioButton || fl.Field().String() == constants.QuestionAttributeOptionTypeCheckbox
	}

	return true
}

func ValidateQuestionAttributesOptionTypeChoices(fl validator.FieldLevel) bool {
	questionType := getQuestionType(fl)
	if questionType == constants.QuestionTypeOption {
		return fl.Field().Len() > 0
	}

	return true
}

func ValidateQuestionAttributesDateTimeType(fl validator.FieldLevel) bool {
	questionType := getQuestionType(fl)
	fmt.Println("questionTypefl", fl.Field().String())
	fmt.Println("questionTypefl", questionType)
	if questionType == constants.QuestionTypeDateTime {
		return fl.Field().String() == constants.QuestionAttributeForematDateOnly || fl.Field().String() == constants.QuestionAttributeFormatDateTime
	}

	return true
}

func ValidateQuestionAttributesInputValidation(fl validator.FieldLevel) bool {
	questionType := getQuestionType(fl)
	if questionType == constants.QuestionTypeShortText {
		return fl.Field().String() == constants.QuestionAttributeInputValidationFreeText || fl.Field().String() == constants.QuestionAttributeInputValidationNumber || fl.Field().String() != constants.QuestionAttributeInputValidationEmail
	}

	return true
}

func getQuestionType(fl validator.FieldLevel) string {
	topVal := reflect.ValueOf(fl.Top().Interface())

	// Walk recursively to find the parent NGOFundingProgramQuestionRequest
	if topVal.Kind() == reflect.Struct {
		qField := topVal.FieldByName("Question")
		if qField.IsValid() && qField.Kind() == reflect.Slice {
			for i := 0; i < qField.Len(); i++ {
				item := qField.Index(i)
				if !item.IsValid() {
					continue
				}

				attrField := item.FieldByName("Attributes")
				if attrField.IsValid() && attrField.CanAddr() {
					if attrField.Addr().Interface() == fl.Parent().Addr().Interface() {
						qType := item.FieldByName("QuestionType")
						if qType.IsValid() && qType.Kind() == reflect.String {
							return qType.String()
						}
					}
				}
			}
		}
	}
	return ""
}

func InitValidator() *validator.Validate {
	validate := validator.New()

	validate.RegisterValidation("date", ValidateDate)
	validate.RegisterValidation("required_by_question_type_text", ValidateQuestionAttributesTextType)
	validate.RegisterValidation("required_by_question_type_option", ValidateQuestionAttributesOptionType)
	validate.RegisterValidation("required_by_question_type_option_choices", ValidateQuestionAttributesOptionTypeChoices)
	validate.RegisterValidation("required_by_question_type_date_time", ValidateQuestionAttributesDateTimeType)
	validate.RegisterValidation("required_if_short_text", ValidateQuestionAttributesInputValidation)

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return validate
}
func mapValidationError(err validator.FieldError) error {
	fieldErrorMap := map[string]map[string]error{
		"title": {
			"required": apperror.ErrNGOFundingProgramTitleRequired,
			"min":      apperror.ErrNGOFundingProgramTitleMinMaxLength,
			"max":      apperror.ErrNGOFundingProgramTitleMinMaxLength,
		},
		"description": {
			"required": apperror.ErrNGOFundingProgramDescriptionRequired,
		},
		"cover_image": {
			"required": apperror.ErrNGOFundingProgramCoverImageRequired,
		},
		"supporting_document": {
			// no validation → no error mapping
		},
		"submission_deadline": {
			"required": apperror.ErrNGOFundingProgramSubmissionDeadlineRequired,
			"date":     apperror.ErrNGOFundingProgramSubmissionDeadlineInvalid,
		},
		"program_period": {
			"required": apperror.ErrNGOFundingProgramPeriodRequired,
		},
		"implementation_cost": {
			"gt": apperror.ErrNGOFundingProgramImplementationCostGT0,
		},
		"question": {
			// optional, no top-level validation
		},
		"question.attributes.placeholder": {
			"required_by_question_type_text": apperror.ErrNGOFundingProgramImplementationCostGT0,
		},
	}

	if fieldErrors, ok := fieldErrorMap[err.Field()]; ok {
		fmt.Println("fieldErrors", err.Field())
		if mappedError, exists := fieldErrors[err.Tag()]; exists {
			fmt.Println("mappedError", err.Tag())
			return mappedError
		}
	}

	return apperror.ErrInvalidRequest

}

func main() {
	request := ProgramCreateRequest{
		Title:              "Judul Program",
		Description:        "Program bantuan sosial",
		CoverImage:         "https://example.com/image.jpg",
		SubmissionDeadline: "2025-12-01",
		ProgramPeriod:      "Oktober - Desember 2025",
		ImplementationCost: 500000000,
		Question: []ProgramQuestionRequest{
			{
				QuestionType: "date_time",
				Label:        "Tanggal Pelaksanaan",
				Attributes: ProgramQuestionAttr{
					Format: "date_time",
				},
				ShowingOrder: 1,
				IsRequired:   true,
			},
			{
				QuestionType: "option",
				Label:        "Jenis Bantuan",
				Attributes: ProgramQuestionAttr{
					OptionType: "radio_button",
					Choices:    []string{"Dana", "Pelatihan"},
				},
				ShowingOrder: 2,
				IsRequired:   true,
			},
		},
	}

	v := InitValidator()

	if err := v.Struct(request); err != nil {
		validationErrors := err.(validator.ValidationErrors)

		for _, fieldErr := range validationErrors {
			fmt.Println("StructField", fieldErr.StructField())
			fmt.Println("fieldErr", fieldErr)
			fmt.Println("Error string", mapValidationError(fieldErr))
		}
	}

	return
}
