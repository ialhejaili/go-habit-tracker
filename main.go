package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/ialhejaili/go-habit-tracker/model"
	"github.com/ialhejaili/go-habit-tracker/repository"
	"github.com/manifoldco/promptui"
)

var (
	db           *sql.DB
	currentUser  *model.User
	habits       []model.Habit
	displayIDMap = make(map[int]int)
	idMap        = make(map[int]int)
)

func main() {

	db = repository.InitDB()
	defer db.Close()

	runCLI()
}

func runCLI() {
	for {
		var items []string

		if currentUser == nil {
			items = []string{"Register", "Login", "Exit"}
		} else {
			items = []string{"Add Habit", "List Habits", "Today's Must Do", "Mark Habit Done", "Delete Habit", "Logout", "Exit"}
		}
		prompt := promptui.Select{
			Label: "Select Action",
			Items: items,
			Size:  10,
		}

		_, result, err := prompt.Run()
		if err != nil {
			log.Fatalf("Prompt failed %v\n", err)
		}

		switch result {
		case "Register":
			registerUser()
			fetchHabits()
		case "Login":
			loginUser()
			fetchHabits()
		case "Add Habit":
			addHabit()
		case "List Habits":
			listHabits()
		case "Delete Habit":
			deleteHabit()
		case "Today's Must Do":
			listHabitsToDoToday()
		case "Mark Habit Done":
			markHabitDone()
		case "Logout":
			logoutUser()
		case "Exit":
			fmt.Println("Goodbye!")
			return
		}
	}
}

func registerUser() {
	username := promptForInput("Enter username")
	password := promptForPassword("Enter password")

	err := repository.RegisterUser(db, username, password)
	if err != nil {
		fmt.Printf("Error registering user: %v\n", err)
		return
	}

	user, err := repository.AuthenticateUser(db, username, password)
	if err != nil {
		fmt.Printf("Error logging in: %v\n", err)
		return
	}
	currentUser = user
	fmt.Printf("User %s registered and logged in successfully!\n", user.Username)
}

func loginUser() {
	username := promptForInput("Enter username")
	password := promptForPassword("Enter password")

	user, err := repository.AuthenticateUser(db, username, password)
	if err != nil {
		fmt.Printf("Error logging in: %v\n", err)
		return
	}
	currentUser = user
	fmt.Printf("User %s logged in successfully!\n", user.Username)
}

func logoutUser() {
	currentUser = nil
	fmt.Println("Logged out successfully!")
}

func promptForInput(label string) string {
	prompt := promptui.Prompt{
		Label: label,
	}
	result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return result
}

func promptForPassword(label string) string {
	prompt := promptui.Prompt{
		Label: label,
		Mask:  '*',
	}
	result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return result
}

func addHabit() {
	if currentUser == nil {
		fmt.Println("Please login first")
		return
	}

	name := promptForInput("Enter habit name")
	description := promptForInput("Enter habit description")

	h := &model.Habit{Name: name, Description: description, UserID: currentUser.ID}
	err := repository.AddHabit(db, h)
	if err != nil {
		fmt.Printf("Error adding habit: %v\n", err)
		return
	}
	fmt.Println("Habit added successfully!")
}

func fetchHabits() {
	if currentUser == nil {
		return
	}
	var err error
	habits, err = repository.ListHabits(db, currentUser.ID)
	if err != nil {
		fmt.Printf("Error fetching habits: %v\n", err)
		return
	}

	if len(habits) > 0 {
		for _, habit := range habits {
			displayIDMap[habit.DisplayID] = habit.ID
			idMap[habit.ID] = habit.DisplayID
		}
	}
}

func listHabits() {
	if currentUser == nil {
		fmt.Println("Please login first")
		return
	}
	var err error
	habits, err = repository.ListHabits(db, currentUser.ID)
	if err != nil {
		fmt.Printf("Error listing habits: %v\n", err)
		return
	}

	if len(habits) > 0 {
		fmt.Println("Your habits:")
		for _, habit := range habits {
			displayIDMap[habit.DisplayID] = habit.ID
			idMap[habit.ID] = habit.DisplayID
			fmt.Printf("%d. %s - %s (Days Completed: %d)\n", habit.DisplayID, habit.Name, habit.Description, habit.DaysCompleted)
		}
	} else {
		fmt.Println("You haven't added any habits yet!")
	}
}

func deleteHabit() {
	if currentUser == nil {
		fmt.Println("Please login first")
		return
	}

	idStr := promptForInput("Enter habit ID to delete")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Printf("Invalid habit ID: %v\n", err)
		return
	}

	err = repository.DeleteHabit(db, currentUser.ID, displayIDMap[id])
	if err != nil {
		fmt.Printf("Error deleting habit: %v\n", err)
		return
	}
	fmt.Println("Habit deleted successfully!")
}

func listHabitsToDoToday() {
	if currentUser == nil {
		fmt.Println("Please login first")
		return
	}

	habits, err := repository.ListHabitsToDoToday(db, currentUser.ID)
	if err != nil {
		fmt.Printf("Error getting habits must do today: %v\n", err)
		return
	}

	if len(habits) == 0 {
		fmt.Println("All habits are done for today!")
		return
	}

	fmt.Println("Habits must do today:")
	for _, habit := range habits {
		habit.DisplayID = idMap[habit.ID]
		fmt.Printf("%d. %s - %s (Days Completed: %d)\n", habit.DisplayID, habit.Name, habit.Description, habit.DaysCompleted)
	}
}

func markHabitDone() {
	if currentUser == nil {
		fmt.Println("Please login first")
		return
	}

	habits, err := repository.ListHabitsToDoToday(db, currentUser.ID)
	if err != nil {
		fmt.Printf("Error getting habits to do today: %v\n", err)
		return
	}

	if len(habits) == 0 {
		fmt.Println("All habits are already done for today!")
		return
	}

	var habitItems []string
	habitMap := make(map[string]int)
	for _, habit := range habits {
		habitItems = append(habitItems, habit.Name)
		habitMap[habit.Name] = habit.ID
	}

	prompt := promptui.Select{
		Label: "Select Habit to Mark Done",
		Items: habitItems,
	}

	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	habitID := habitMap[result]
	selectedHabit := model.GetHabitByID(habits, habitID)
	err = repository.UpdateHabitDone(db, selectedHabit.ID, selectedHabit.DaysCompleted)
	if err != nil {
		fmt.Printf("Error marking habit as done: %v\n", err)
		return
	}

	fmt.Printf("Habit '%s' marked as done for today!\n", result)
}
