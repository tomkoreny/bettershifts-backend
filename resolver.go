//go:generate go run scripts/gqlgen.go -v
package bettershifts

import (
	"context"
	"errors"
  "fmt" 
	"crypto/rand"
  "time"
	"github.com/jinzhu/gorm"
  "github.com/google/uuid"

	"github.com/lordpuma/bettershifts/models"
	"github.com/lordpuma/bettershifts/auth"
)

type Resolver struct{
	Db *gorm.DB
}

func (r *Resolver) Benefit() BenefitResolver {
	return &benefitResolver{r}
}
func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}
func (r *Resolver) Shift() ShiftResolver {
	return &shiftResolver{r}
}
func (r *Resolver) Todo() TodoResolver {
	return &todoResolver{r}
}
func (r *Resolver) Workplace() WorkplaceResolver {
	return &workplaceResolver{r}
}

type benefitResolver struct{ *Resolver }

func (r *benefitResolver) User(ctx context.Context, obj *models.Benefit) (models.User, error) {
	var user models.User
	r.Db.First(&user, "id", obj.UserID)
	return user, nil
}
func (r *benefitResolver) Date(ctx context.Context, obj *models.Benefit) (string, error) {
  return obj.Date.Format(time.RFC3339), nil
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateUser(ctx context.Context, input NewUser) (models.User, error) {
  if user := auth.ForContext(ctx) ; user == nil || !user.IsAdmin {
		return models.User{}, fmt.Errorf("Access denied")
	}

	user := models.User{
    ID: uuid.New().String(),
		FirstName: input.FirstName,
		LastName: input.LastName,
		UserName: input.Username,
		IsAdmin: input.IsAdmin,
		Wage: input.Wage,
	}
	r.Db.Create(&user)
	return user, nil
}
func (r *mutationResolver) CreateTodo(ctx context.Context, input NewTodo) (models.Todo, error) {
  if user := auth.ForContext(ctx) ; user == nil || !user.IsAdmin {
		return models.Todo{}, fmt.Errorf("Access denied")
	}
	todo := models.Todo{
    ID: uuid.New().String(),
		Name:   input.Name,
		WorkplaceID: input.Workplace,
		Benefit: input.Benefit,
	}
  newdate, _ := time.Parse(time.RFC3339, input.Date)
  todo.Date = newdate
	r.Db.Create(&todo)
	return todo, nil
}
func (r *mutationResolver) CreateShift(ctx context.Context, input NewShift) (models.Shift, error) {
  if user := auth.ForContext(ctx) ; user == nil || !user.IsAdmin {
		return models.Shift{}, fmt.Errorf("Access denied")
	}
	shift := models.Shift{
    ID: uuid.New().String(),
		WorkplaceID: input.Workplace,
		UserID: input.User,
	}
	shift.Start, _ = time.Parse(time.RFC3339, input.Start)
  if input.End != nil {
	  *shift.End, _ = time.Parse(time.RFC3339, *input.End) 
  }
	r.Db.Create(&shift)
	return shift, nil
}
func (r *mutationResolver) CreateWorkplace(ctx context.Context, input NewWorkplace) (models.Workplace, error) {
  if user := auth.ForContext(ctx) ; user == nil || !user.IsAdmin {
		return models.Workplace{}, fmt.Errorf("Access denied")
	}
	workplace := models.Workplace{
    ID: uuid.New().String(),
		Name: input.Name,
	}
	r.Db.Create(&workplace)
	return workplace, nil
}
func (r *mutationResolver) CreateBenefit(ctx context.Context, input NewBenefit) (models.Benefit, error) {
  if user := auth.ForContext(ctx) ; user == nil || !user.IsAdmin {
		return models.Benefit{}, fmt.Errorf("Access denied")
	}
	benefit := models.Benefit{
      ID: uuid.New().String(),
     	Reason: input.Reason,
     	Amount: input.Amount,
     	UserID: input.User,
	}
  newdate, _ := time.Parse(time.RFC3339, input.Date)
  benefit.Date = newdate
	r.Db.Create(&benefit)
	return benefit, nil
}
func (r *mutationResolver) EditUser(ctx context.Context, id string, input EditUser) (models.User, error) {
  if user := auth.ForContext(ctx) ; user == nil || !user.IsAdmin {
		return models.User{}, fmt.Errorf("Access denied")
	}
	var user models.User
  r.Db.First(&user, models.User{ID: id})
	if user.ID == "" {
		return user, errors.New("user-not-found")
	}
	if input.FirstName != nil {
		user.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		user.LastName = *input.LastName
	}
	if input.Username != nil {
		user.UserName = *input.Username
	}
	if input.IsAdmin != nil {
		user.IsAdmin = *input.IsAdmin
	}
	if input.Wage != nil {
		user.Wage = *input.Wage
	}
	if input.Workplaces != nil {
		var workplaces []models.Workplace 
		r.Db.Where("id in (?)", input.Workplaces).Find(&workplaces)
		r.Db.Model(&user).Association("Workplaces").Replace(workplaces) 
	}
	r.Db.Model(&user).Where("id", user.ID).Omit("workplaces").Update(user)
	return user, nil
}
func (r *mutationResolver) EditTodo(ctx context.Context, id string, input EditTodo) (models.Todo, error) {
  if user := auth.ForContext(ctx) ; user == nil || !user.IsAdmin {
		return models.Todo{}, fmt.Errorf("Access denied")
	}
	var todo models.Todo
	r.Db.First(&todo, "id", id)
	if todo.ID == "" {
		return todo, errors.New("todo-not-found")
	}
	if input.Date != nil {
    newdate, _ := time.Parse(time.RFC3339, *input.Date)
    todo.Date = newdate
	}
	if input.Benefit != nil {
		todo.Benefit= *input.Benefit
	}
	if input.Name != nil {
		todo.Name= *input.Name
	}
	if input.Workplace != nil {
    var workplace models.Workplace
    r.Db.First(&workplace, "id", input.Workplace)
    if workplace.ID == "" {
      return todo, errors.New("workplace-does-not-exist")
    }
    todo.WorkplaceID = workplace.ID
	}
	r.Db.Model(&todo).Updates(todo)
	return todo, nil
}
func (r *mutationResolver) EditShift(ctx context.Context, id string, input EditShift) (models.Shift, error) {
  if user := auth.ForContext(ctx) ; user == nil || !user.IsAdmin {
		return models.Shift{}, fmt.Errorf("Access denied")
	}
  var shift models.Shift
  r.Db.First(&shift, "id", id)
  if shift.ID == "" {
    return shift, errors.New("shift-not-found")
  }
	if input.Start != nil {
		shift.Start, _ = time.Parse(time.RFC3339, *input.Start)
	}
	if input.End != nil {
		*shift.End, _ = time.Parse(time.RFC3339, *input.End)
	}
	if input.Workplace != nil {
    var workplace models.Workplace
    r.Db.First(&workplace, "id", input.Workplace)
    if workplace.ID == "" {
      return shift, errors.New("workplace-does-not-exist")
    }
    shift.WorkplaceID = workplace.ID
	}
	if input.User != nil {
    var user models.User
    r.Db.First(&user, "id", input.User)
    if user.ID == "" {
      return shift, errors.New("user-does-not-exist")
    }
    shift.UserID = user.ID
	}
	r.Db.Model(&shift).Updates(shift)
  return shift, nil
} 
func (r *mutationResolver) EditWorkplace(ctx context.Context, id string, input EditWorkplace) (models.Workplace, error) {
  if user := auth.ForContext(ctx) ; user == nil || !user.IsAdmin {
		return models.Workplace{}, fmt.Errorf("Access denied")
	}
  var workplace models.Workplace
  r.Db.First(&workplace, "id", id)
  if workplace.ID == "" {
    return workplace, errors.New("workplace-not-found")
  }
  if input.Name != nil {
    workplace.Name = *input.Name
  }
	r.Db.Model(&workplace).Updates(workplace)
  return workplace, nil
}
func (r *mutationResolver) EditBenefit(ctx context.Context, id string, input EditBenefit) (models.Benefit, error) {
  if user := auth.ForContext(ctx) ; user == nil || !user.IsAdmin {
		return models.Benefit{}, fmt.Errorf("Access denied")
	}
	var benefit models.Benefit
	r.Db.First(&benefit, "id", id)
	if benefit.ID == "" {
	  return benefit, errors.New("benefit-not-found")
	}
	if input.User != nil { 
    var user models.User
    r.Db.First(&user, "id", input.User)
    if user.ID == "" {
      return benefit, errors.New("user-does-not-exist")
    }
    benefit.UserID = user.ID
	}
  if input.Amount != nil {
    benefit.Amount = *input.Amount
  }
  if input.Date != nil {
    newdate, _ := time.Parse(time.RFC3339, *input.Date)
    benefit.Date = newdate
  }
  if input.Reason != nil {
    benefit.Reason = *input.Reason
  }
	r.Db.Model(&benefit).Updates(benefit)
  return benefit, nil
}
func (r *mutationResolver) ChangePassword(ctx context.Context, oldpass string, newpass string) (bool, error) {
  user := auth.ForContext(ctx)
  if user == nil {
		return false, fmt.Errorf("User not logged in")
	}
  if user.Password == "" || auth.ComparePasswords(user.Password, oldpass) {
    r.Db.Model(&user).Update("Password", auth.HashAndSalt(newpass))
    return true, nil
  }
return false, fmt.Errorf("Wrong old pass")
}
func (r *mutationResolver) Login(ctx context.Context, username string, password string) (LoginPayload, error) { 
  var user models.User
  r.Db.First(&user, models.User{UserName: username})
  if user.ID != "" && (user.Password == "" || auth.ComparePasswords(user.Password, password)) {
    var token models.Token
    tokenstring, _ := GenerateRandomString(32)
    token = models.Token{
      ID: uuid.New().String(),
      UserID: user.ID,
      Token: tokenstring,
    }
    r.Db.Save(&token)
    return LoginPayload{ Token: token.Token, HasPassword: user.Password == "" }, nil
  }
  return LoginPayload{}, errors.New("Invalid username or password")
}
func (r *mutationResolver) StartShift(ctx context.Context, workplace string) (bool, error) {
  user := auth.ForContext(ctx)
  if user == nil {
		return false, fmt.Errorf("User not logged in")
	}
  var workplaceModel models.Workplace
  r.Db.First(&workplaceModel, models.Workplace{ID: workplace})
  if workplaceModel.ID == "" {
    return false, errors.New("workplace-not-found")
  }
	shift := models.Shift{
    ID: uuid.New().String(),
		WorkplaceID: workplaceModel.ID,
		UserID: user.ID,
    Start: time.Now(),
	}
	r.Db.Create(&shift)
	return true, nil
}
func (r *mutationResolver) StopShift(ctx context.Context, finishedTodos []string) (bool, error) {
  user := auth.ForContext(ctx)
  if user == nil {
		return false, fmt.Errorf("User not logged in")
	}
  var shift models.Shift
  r.Db.Where("user_id = ? AND end is null", user.ID).Order("start desc").First(&shift)
  if shift.ID == "" {
		return false, fmt.Errorf("User does not have a shift")
  }
  time := time.Now()
  r.Db.Model(&shift).Where(shift).Update("end", &time)
  return true, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Users(ctx context.Context) ([]models.User, error) {
  if user := auth.ForContext(ctx) ; user == nil || !user.IsAdmin {
		return []models.User{}, fmt.Errorf("Access denied")
	}
	var users []models.User
	r.Db.Find(&users)
	return users, nil
}
func (r *queryResolver) UserByID(ctx context.Context, id string) (models.User, error) {
  if user := auth.ForContext(ctx) ; user == nil || !user.IsAdmin {
		return models.User{}, fmt.Errorf("Access denied")
	}
	var user models.User
	r.Db.First(&user, "id", id)
  if user.ID == "" {
    return user, errors.New("user-not-found")
  }
	return user, nil
}
func (r *queryResolver) Workplaces(ctx context.Context) ([]models.Workplace, error) {
  if user := auth.ForContext(ctx) ; user == nil || !user.IsAdmin {
		return []models.Workplace{}, fmt.Errorf("Access denied")
	}
	var workplaces []models.Workplace
	r.Db.Find(&workplaces)
	return workplaces, nil
}
func (r *queryResolver) WorkplaceByID(ctx context.Context, id string) (models.Workplace, error) {
  if user := auth.ForContext(ctx) ; user == nil || !user.IsAdmin {
		return models.Workplace{}, fmt.Errorf("Access denied")
	}
	var workplace models.Workplace
	r.Db.First(&workplace, "id", id)
  if workplace.ID == "" {
    return workplace, errors.New("workplace-not-found")
  }
	return workplace, nil
}
func (r *queryResolver) Todos(ctx context.Context, date *string, workplace *string) ([]models.Todo, error) {
  if user := auth.ForContext(ctx) ; user == nil || !user.IsAdmin {
		return []models.Todo{}, fmt.Errorf("Access denied")
	}
	var todos []models.Todo
	r.Db.Find(&todos)
	return todos, nil
}
func (r *queryResolver) TodoByID(ctx context.Context, id string) (models.Todo, error) {
  if user := auth.ForContext(ctx) ; user == nil || !user.IsAdmin {
		return models.Todo{}, fmt.Errorf("Access denied")
	}
	var todo models.Todo
	r.Db.First(&todo, "id", id)
  if todo.ID == "" {
    return todo, errors.New("todo-not-found")
  }
	return todo, nil
}
func (r *queryResolver) Shifts(ctx context.Context, since *string, untill *string, workplace *string, user *string) ([]models.Shift, error) {
  if user := auth.ForContext(ctx) ; user == nil || !user.IsAdmin {
		return []models.Shift{}, fmt.Errorf("Access denied")
	}
	var shifts []models.Shift
	r.Db.Find(&shifts)
	return shifts, nil
}
func (r *queryResolver) ShiftByID(ctx context.Context, id string) (models.Shift, error) {
  if user := auth.ForContext(ctx) ; user == nil || !user.IsAdmin {
		return models.Shift{}, fmt.Errorf("Access denied")
	}
	var shift models.Shift
	r.Db.First(&shift, "id", id)
  if shift.ID == "" {
    return shift, errors.New("shift-not-found")
  }
	return shift, nil
}
func (r *queryResolver) Benefits(ctx context.Context, since *string, untill *string, workplace *string, user *string) ([]models.Benefit, error) {
  if user := auth.ForContext(ctx) ; user == nil || !user.IsAdmin {
		return []models.Benefit{}, fmt.Errorf("Access denied")
	}
	var benefits []models.Benefit
	r.Db.Find(&benefits)
	return benefits, nil
}
func (r *queryResolver) BenefitByID(ctx context.Context, id string) (models.Benefit, error) {
  if user := auth.ForContext(ctx) ; user == nil || !user.IsAdmin {
		return models.Benefit{}, fmt.Errorf("Access denied")
	}
	var benefit models.Benefit
	r.Db.First(&benefit, "id", id)
  if benefit.ID == "" {
    return benefit, errors.New("benefit-not-found")
  }
	return benefit, nil
}
func (r *queryResolver) Wages(ctx context.Context, month *string) ([]Wage, error) {
  user := auth.ForContext(ctx) 
  if user == nil || !user.IsAdmin {
		return []Wage{}, fmt.Errorf("Access denied")
	}
  var duration []int
	r.Db.Raw("select DATEDIFF('mi', start, end) as duration from shifts where start >= '?' AND start <= '?'", "2019-02-01", "2019-02-28").Pluck("duration", &duration)
  fmt.Println(duration)
	panic("not implemented")
}

type shiftResolver struct{ *Resolver }

func (r *shiftResolver) Workplace(ctx context.Context, obj *models.Shift) (models.Workplace, error) {
	var workplace models.Workplace
	r.Db.First(&workplace, "id", obj.WorkplaceID)
	return workplace, nil
}
func (r *shiftResolver) User(ctx context.Context, obj *models.Shift) (models.User, error) {
	var user models.User
	r.Db.First(&user, "id", obj.UserID)
	return user, nil
}
func (r *shiftResolver) Start(ctx context.Context, obj *models.Shift) (string, error) {
  return obj.Start.Format(time.RFC3339), nil
}
func (r *shiftResolver) End(ctx context.Context, obj *models.Shift) (string, error) {
  return obj.End.Format(time.RFC3339), nil
}

type todoResolver struct{ *Resolver }

func (r *todoResolver) DoneBy(ctx context.Context, obj *models.Todo) (models.User, error) {
	var user models.User
	r.Db.First(&user, "id", obj.UserID)
	return user, nil
}
func (r *todoResolver) Workplace(ctx context.Context, obj *models.Todo) (models.Workplace, error) {
	var workplace models.Workplace
	r.Db.First(&workplace, "id", obj.WorkplaceID)
	return workplace, nil
}
func (r *todoResolver) Date(ctx context.Context, obj *models.Todo) (string, error) {
  return obj.Date.Format(time.RFC3339), nil
}

type workplaceResolver struct{ *Resolver }

func (r *workplaceResolver) Users(ctx context.Context, obj *models.Workplace) ([]models.User, error) {
	var user []models.User
	r.Db.Find(&user)
	return user, nil
}

func GenerateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	bytes, err := GenerateRandomBytes(n)
	if err != nil {
		return "", err
	}
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes), nil
}

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

