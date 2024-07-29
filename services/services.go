package services

type ServicesInterface interface {
	AuthServiceInterface
	UserServiceInterface
}

type services struct {
	AuthServiceInterface
	UserServiceInterface
}

func NewServices(authService AuthServiceInterface, userService UserServiceInterface) ServicesInterface {
	return &services{
		AuthServiceInterface: authService,
		UserServiceInterface: userService,
	}
}
