package main

import (
	"fmt"
	"log"
	"net/mail"
	"os"
	"strings"
	"unicode"

	"github.com/jasonuc/usermakertui/db"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"golang.org/x/crypto/bcrypt"
)

var logOut = log.New(os.Stdout, "usermaker: ", log.Lmsgprefix+log.Ldate+log.Ltime)
var logErr = log.New(os.Stderr, "usermaker: [ERROR] ", log.Lmsgprefix+log.Ldate+log.Ltime)

type adduserForm struct {
	inputs []textinput.Model
	focus  int
	state  formState
}

type move int
type formState int
type overflow bool
type listErrors struct {
	errored []bool
	items   []string
}

func (l listErrors) Error() string {
	var msg strings.Builder

	for idx, item := range l.items {
		if idx > 0 {
			msg.WriteRune('\n')
		}
		if l.errored[idx] {
			msg.WriteString(redCross)
		} else {
			msg.WriteString(greenCheck)
		}
		msg.WriteRune(' ')
		msg.WriteString(item)
	}

	return msg.String()
}

const (
	adduserEmail = iota
	adduserPwd
	filling formState = iota
	saving
	cancelling
	down   move     = 1
	up     move     = -1
	rotate overflow = true
	noop   overflow = false
)

var (
	greenCheck = lipgloss.NewStyle().Foreground(lipgloss.Color("70")).Render("✓")
	redCross   = lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Render("✗")
)

var (
	adduserTitle string = lipgloss.NewStyle().
			Bold(true).
			Background(lipgloss.Color("63")).
			Foreground(lipgloss.Color("228")).
			Render(" New User ")
	submitMsg string = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Render("press Enter until the end or Ctrl+s to submit the form")
)

func main() {
	db.InitMockDB() // Initialize the mock database
	cli := tea.NewProgram(initialModel())
	model, formErr := cli.Run()
	if formErr != nil {
		logErr.Fatal(formErr)
	}
	form, _ := model.(adduserForm)
	if form.state != saving {
		logOut.Print("You cancelled form submission of new user")
		os.Exit(0)
	}
	email := form.inputs[adduserEmail].Value()
	pwd := form.inputs[adduserPwd].Value()

	validEmail, mailErr := validateEmail(email)
	if mailErr != nil {
		logErr.Fatalf("Email '%s' is invalid: %s", email, strings.ReplaceAll(mailErr.Error(), "\n", " "))
	}

	hash, pwdErr := validatePassword(pwd)
	if pwdErr != nil {
		logErr.Fatalf("Password is invalid: %s", strings.ReplaceAll(pwdErr.Error(), "\n", " "))
	}

	logOut.Printf("Creating user <%s>", validEmail)
	user, createErr := db.Q.CreateUser(db.CreateUserParams{
		Email:    validEmail,
		Password: hash,
	})
	if createErr != nil {
		logErr.Fatal(createErr)
	}

	logOut.Printf("Created user #%d <%s>", user.ID, validEmail)
}

func initialModel() adduserForm {
	emailIn := textinput.New()
	emailIn.Prompt = ""
	emailIn.Placeholder = "firstname.surname@domain.com"
	emailIn.Focus()
	emailIn.Validate = func(s string) error {
		_, err := validateEmail(s)
		return err
	}

	pwdIn := textinput.New()
	pwdIn.Prompt = ""
	pwdIn.EchoMode = textinput.EchoPassword
	pwdIn.Placeholder = "v3rySTRONG-p@55word (don't copy this one…)"
	pwdIn.Validate = func(s string) error {
		_, err := validatePassword(s)
		return err
	}

	return adduserForm{
		inputs: []textinput.Model{emailIn, pwdIn},
		focus:  adduserEmail,
		state:  filling,
	}
}

func (form adduserForm) Init() tea.Cmd {
	return textinput.Blink
}

func (form adduserForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			form.state = cancelling
			return form, tea.Quit
		case tea.KeyDown:
			if form.changeFocus(down, noop) {
				return form, textinput.Blink
			}
		case tea.KeyUp:
			if form.changeFocus(up, noop) {
				return form, textinput.Blink
			}
		case tea.KeyCtrlS, tea.KeyEnter:
			if form.changeFocus(down, noop) {
				return form, textinput.Blink
			} else {
				form.state = saving
				return form, tea.Quit
			}
		case tea.KeyTab:
			form.changeFocus(down, rotate)
			return form, textinput.Blink
		case tea.KeyShiftTab:
			form.changeFocus(up, rotate)
			return form, textinput.Blink
		}
	}

	var cmd tea.Cmd
	form.inputs[form.focus], cmd = form.inputs[form.focus].Update(msg)
	return form, cmd
}

func (form *adduserForm) changeFocus(move move, overflow overflow) bool {
	dest := form.focus + int(move)
	if overflow || 0 <= dest && dest < len(form.inputs) {
		form.inputs[form.focus].Blur()
		form.focus = (len(form.inputs) + form.focus + int(move)) % len(form.inputs)
		form.inputs[form.focus].Focus()
		return true
	}
	return false
}

func (form adduserForm) View() string {
	if form.state != filling {
		return ""
	}

	var emailCell strings.Builder
	emailCell.WriteString(form.inputs[adduserEmail].View())
	emailCell.WriteRune('\n')
	if form.inputs[adduserEmail].Value() == "" {
		emailCell.WriteString(redCross)
		emailCell.WriteString(" required")
	} else if form.inputs[adduserEmail].Err == nil {
		emailCell.WriteString(greenCheck)
		emailCell.WriteString(" good")
	} else {
		emailCell.WriteString(form.inputs[adduserEmail].Err.Error())
	}

	var pwdCell strings.Builder
	pwdCell.WriteString(form.inputs[adduserPwd].View())
	pwdCell.WriteRune('\n')
	if form.inputs[adduserPwd].Value() == "" {
		pwdCell.WriteString(redCross)
		pwdCell.WriteString(" required")
	} else if form.inputs[adduserPwd].Err == nil {
		pwdCell.WriteString(greenCheck)
		pwdCell.WriteString(" good")
	} else {
		pwdCell.WriteString(form.inputs[adduserPwd].Err.Error())
	}

	cellStyle := lipgloss.NewStyle().Padding(0, 1)
	t := table.New().
		Border(lipgloss.HiddenBorder()).
		BorderRow(true).
		Rows(
			[]string{"email", emailCell.String()},
			[]string{"password", pwdCell.String()},
		).
		StyleFunc(func(row, col int) lipgloss.Style {
			if col == 0 {
				return cellStyle.Foreground(lipgloss.Color("63"))
			}
			return cellStyle
		})

	return lipgloss.NewStyle().Padding(1, 2).Render(fmt.Sprintf("%s\n%s\n%s",
		adduserTitle,
		t.Render(),
		submitMsg,
	))
}

func validateEmail(email string) (string, error) {
	err := emailNoError()

	mailAddr, mailErr := mail.ParseAddress(email)
	if mailErr != nil {
		err.errored[0] = true
		return "", err
	}

	user, _ := db.Q.SearchUser(mailAddr.Address)
	if user.ID > 0 {
		err.errored[1] = true
		return "", err
	}

	return mailAddr.Address, nil
}

func validatePassword(cleartext string) (string, error) {
	err := pwdNoError()

	err.errored[0] = len(cleartext) < 10
	hash, bcryptErr := bcrypt.GenerateFromPassword([]byte(cleartext), bcrypt.DefaultCost)
	if bcryptErr != nil {
		err.errored[1] = true
	}

	err.errored[2] = true
	err.errored[3] = true
	err.errored[4] = true
	for _, c := range cleartext {
		switch {
		case unicode.IsLower(c):
			err.errored[2] = false
		case unicode.IsUpper(c):
			err.errored[3] = false
		default:
			err.errored[4] = false
		}
	}

	if Any(err.errored) {
		return "", err
	} else {
		return string(hash), nil
	}
}

func emailNoError() listErrors {
	return listErrors{
		errored: []bool{false, false},
		items:   []string{"is a valid email address", "is available"},
	}
}

func pwdNoError() listErrors {
	return listErrors{
		errored: []bool{false, false, false, false, false},
		items: []string{
			"has 10 characters or more",
			"has less than 72 bytes",
			"has a lowercase letter",
			"has an uppercase letter",
			"has a digit or special character",
		},
	}
}

func Any(bools []bool) bool {
	for _, b := range bools {
		if b {
			return true
		}
	}
	return false
}
