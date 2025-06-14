package graphQlDemo

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.72

import (
	"context"
	"fmt"
	"graphQlDemo/auth"
	"graphQlDemo/ent"
	"graphQlDemo/ent/review"
	"graphQlDemo/ent/tool"
	"graphQlDemo/ent/user"
	"graphQlDemo/utils"

	"github.com/google/uuid"
	"github.com/nyaruka/phonenumbers"
)

// Createreview is the resolver for the createreview field.
func (r *mutationResolver) Createreview(ctx context.Context, input ent.CreateReviewInput) (*ent.Review, error) {
	usr, ok := auth.UserFromContext(ctx)
	if !ok || usr.Role != "user" {
		return nil, fmt.Errorf("access denied")
	}

	if input.ReviwedToolID == nil {
		return nil, fmt.Errorf("reviwedToolID is required")
	}

	exists, existsErr := r.client.Review.
		Query().
		Where(
			review.HasReviewerWith(user.ID(usr.ID)),
			review.HasReviwedToolWith(tool.ID(*input.ReviwedToolID)),
		).
		Exist(ctx)
	if existsErr != nil {
		return nil, fmt.Errorf("failed to check existing reviews")
	}

	if exists {
		return nil, fmt.Errorf("you have already reviewed this tool")
	}

	userUUID, userUUIDErr := uuid.Parse(usr.ID.String())
	if userUUIDErr != nil {
		return nil, fmt.Errorf("invalid user ID")
	}

	input.ReviewerID = &userUUID
	review, reviewErr := r.client.Review.Create().
		SetRating(input.Rating).
		SetComment(input.Comment).
		SetReviewerID(*input.ReviewerID).
		SetReviwedToolID(*input.ReviwedToolID).
		Save(ctx)
	if reviewErr != nil {
		return nil, fmt.Errorf("failed to create review")
	}

	updateErr := utils.UpdateToolRating(ctx, r.client, *input.ReviwedToolID)
	if updateErr != nil {
		return nil, fmt.Errorf("failed to update tool rating")
	}

	return review, nil
}

// Createuser is the resolver for the createuser field.
func (r *mutationResolver) Createuser(ctx context.Context, input ent.CreateUserInput) (*ent.User, error) {
	currentUser, ok := auth.UserFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("no user in context")
	}

	if currentUser.Role != "admin" {
		return nil, fmt.Errorf("admin role is required")
	}

	hashedPass, hashedPassErr := auth.HashPassword(input.PasswordHash)
	if hashedPassErr != nil {
		return nil, fmt.Errorf("error in create user")
	}
	return r.client.User.Create().
		SetName(input.Name).
		SetUsername(input.Username).
		SetEmail(input.Email).
		SetPasswordHash(hashedPass).
		Save(ctx)
}

// Createtool is the resolver for the createtool field.
func (r *mutationResolver) Createtool(ctx context.Context, input ent.CreateToolInput) (*ent.Tool, error) {
	currentUser, ok := auth.UserFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("no user in context")
	}

	if currentUser.Role != "admin" {
		return nil, fmt.Errorf("admin role is required")
	}

	return r.client.Tool.Create().
		SetName(input.Name).
		SetDescription(input.Description).
		SetCategory(input.Category).
		SetWebsite(input.Website).
		SetImageURL(input.ImageURL).
		Save(ctx)
}

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, username string, password string) (*string, error) {
	user, userErr := r.client.User.
		Query().
		Where(user.Username(username)).
		Only(ctx)

	if userErr != nil {
		return nil, fmt.Errorf("password or user invalid")
	}
	if !auth.CheckPassword(password, user.PasswordHash) {
		return nil, fmt.Errorf("password or user invalid")
	}
	token, tokenErr := auth.CreateToken(user.Username)
	if tokenErr != nil {
		return nil, fmt.Errorf("password or user invalid")
	}

	return &token, tokenErr
}

// Register is the resolver for the register field.
func (r *mutationResolver) Register(ctx context.Context, username string, password string, email string) (*string, error) {
	exists, existsErr := r.client.User.
		Query().
		Where(user.Or(
			user.Username(username),
			user.Email(email),
		)).
		Exist(ctx)
	if existsErr != nil {
		return nil, fmt.Errorf("error checking if user exist")
	}
	if exists {
		return nil, fmt.Errorf("username or email already in use")
	}

	hashedPass, hashedPassErr := auth.HashPassword(password)
	if hashedPassErr != nil {
		return nil, fmt.Errorf("error hashing password")
	}

	_, createErr := r.client.User.Create().
		SetUsername(username).
		SetEmail(email).
		SetPasswordHash(hashedPass).
		SetRole("user").
		SetName(username).
		Save(ctx)
	if createErr != nil {
		return nil, fmt.Errorf("error creating user")
	}

	token, tokenErr := auth.CreateToken(username)
	if tokenErr != nil {
		return nil, fmt.Errorf("error generating token")
	}

	return &token, nil
}

// ChangePassword is the resolver for the changePassword field.
func (r *mutationResolver) ChangePassword(ctx context.Context, currentPassword string, newPassword string) (bool, error) {
	currentUser, ok := auth.UserFromContext(ctx)
	if !ok {
		return false, fmt.Errorf("no user in context")
	}

	valid := auth.CheckPassword(currentPassword, currentUser.PasswordHash)
	if !valid {
		return false, fmt.Errorf("current password is incorrect")
	}

	hashedPass, hashedPassErr := auth.HashPassword(newPassword)
	if hashedPassErr != nil {
		return false, fmt.Errorf("error in hash password")
	}

	r.client.User.UpdateOne(currentUser).SetPasswordHash(hashedPass).Exec(ctx)
	return true, nil
}

// UpdateProfile is the resolver for the updateProfile field.
func (r *mutationResolver) UpdateProfile(ctx context.Context, input UpdateProfileInput) (*ent.User, error) {
	currentUser, ok := auth.UserFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("authentication required")
	}

	update := r.client.User.UpdateOneID(currentUser.ID)
	if input.Name != nil {
		update.SetName(*input.Name)
	}

	if input.Email != nil {
		emailExists, emailExistsErr := r.client.User.
			Query().
			Where(user.EmailEQ(*input.Email)).
			Where(user.IDNEQ(currentUser.ID)).
			Exist(ctx)
		if emailExistsErr != nil {
			return nil, fmt.Errorf("error checking email availability")
		}
		if emailExists {
			return nil, fmt.Errorf("email already in use")
		}
		update.SetEmail(*input.Email)
	}

	if input.Username != nil {
		userNameExists, userNameExistsErr := r.client.User.
			Query().
			Where(user.UsernameEQ(*input.Username)).
			Where(user.IDNEQ(currentUser.ID)).
			Exist(ctx)
		if userNameExistsErr != nil {
			return nil, fmt.Errorf("error checking username availability")
		}
		if userNameExists {
			return nil, fmt.Errorf("username already in use")
		}
		update.SetUsername(*input.Username)
	}

	return update.Save(ctx)
}

// VerifyUser is the resolver for the verifyUser field.
func (r *mutationResolver) VerifyUser(ctx context.Context, phoneNumber string, twilioCode string) (*VerifyUserResponse, error) {
	parsedNum, err := phonenumbers.Parse(phoneNumber, "")
    if err != nil {
        return &VerifyUserResponse{
            Success: false,
            Message: "Invalid phone number format",
            Code:    400,
        }, nil
    }
    formattedNum := phonenumbers.Format(parsedNum, phonenumbers.E164)

    verification, err := utils.VerifyNumber(formattedNum, twilioCode)
    if err != nil {
        return &VerifyUserResponse{
            Success: false,
            Message: "Unable to verify code",
            Code:    500,
        }, nil
    }

    switch *verification {
    case "approved":
        currentUser, ok := auth.UserFromContext(ctx)
        if !ok {
            return &VerifyUserResponse{
                Success: false,
                Message: "Authentication required",
                Code:    401,
            }, nil
        }

        if currentUser.PhoneNumber != formattedNum {
            return &VerifyUserResponse{
                Success: false,
                Message: "Phone number does not match user record",
                Code:    400,
            }, nil
        }

        _, err = r.client.User.UpdateOneID(currentUser.ID).
            SetIsVerified(true).
            Save(ctx)
        if err != nil {
            return &VerifyUserResponse{
                Success: false,
                Message: "Failed to update user verification status",
                Code:    500,
            }, nil
        }

        return &VerifyUserResponse{
            Success: true,
            Message: "Phone number successfully verified",
            Code:    200,
        }, nil

    default:
        return &VerifyUserResponse{
            Success: false,
            Message: "Verification failed",
            Code:    400,
        }, nil
    }
}

// SendVerificationCode is the resolver for the sendVerificationCode field.
func (r *mutationResolver) SendVerificationCode(ctx context.Context, phoneNumber string) (bool, error) {
	currentUser, ok := auth.UserFromContext(ctx)
    if !ok {
        return false, fmt.Errorf("authentication required")
    }

    parsedNum, err := phonenumbers.Parse(phoneNumber, "")
    if err != nil {
        return false, fmt.Errorf("invalid phone number format")
    }
    formattedNum := phonenumbers.Format(parsedNum, phonenumbers.E164)

    _, err = r.client.User.UpdateOneID(currentUser.ID).
        SetPhoneNumber(formattedNum).
        Save(ctx)
    if err != nil {
        return false, fmt.Errorf("failed to update user phone number")
    }

    if err := utils.SendVerificationCode(formattedNum); err != nil {
        return false, err
    }

    return true, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
