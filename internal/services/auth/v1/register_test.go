package authv1

// func (s *UnitTestingSuite) TestRegister(t *testing.T) {
// 	ctx := context.Background()
// 	assert := assert.New(t)
// 	require := require.New(t)
//
// 	s.DB.EXPECT().CreateUser(mock.Anything, mock.Anything).RunAndReturn(
// 		func(_ context.Context, params repository.CreateUserParams) (repository.User, error) {
// 			return repository.User{
// 				ID:        params.ID,
// 				Email:     params.Email,
// 				Password:  params.Password,
// 				Scope:     params.Scope,
// 				CreatedAt: time.Now().UTC(),
// 				UpdatedAt: time.Now().UTC(),
// 			}, nil
// 		})
//
// 	res, err := s.Client.Register(ctx, &authv1.RegisterRequest{
// 		Email:    "test@mail.com",
// 		Password: "1234",
// 		UserType: authv1.UserType_USER_TYPE_APPLICANT,
// 	})
// 	require.NoError(err)
//
// 	assert.NotZero(res.AccessToken)
// 	assert.NotZero(res.RefreshToken)
// }
