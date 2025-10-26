package utils_test

import (
    "testing"
    "toolcat/config"
    "toolcat/utils"
)

func TestGenerateAndVerifyAccessToken(t *testing.T) {
    // Arrange
    config.Config.JWT.Secret = "testsecret"
    config.Config.JWT.AccessTokenExpiry = 60

    userID := uint(123)
    tenantID := uint(456)

    // Act
    token, err := utils.GenerateToken(userID, tenantID)
    if err != nil {
        t.Fatalf("GenerateToken error: %v", err)
    }
    if token == "" {
        t.Fatalf("GenerateToken returned empty token")
    }

    gotUserID, tokenType, gotTenantID, err := utils.VerifyToken(token)
    if err != nil {
        t.Fatalf("VerifyToken error: %v", err)
    }

    // Assert
    if tokenType != "access" {
        t.Errorf("expected token type 'access', got %s", tokenType)
    }
    if gotUserID != userID {
        t.Errorf("expected userID %d, got %d", userID, gotUserID)
    }
    if gotTenantID != tenantID {
        t.Errorf("expected tenantID %d, got %d", tenantID, gotTenantID)
    }
}

func TestGenerateAndVerifyRefreshToken(t *testing.T) {
    // Arrange
    config.Config.JWT.Secret = "testsecret"
    config.Config.JWT.RefreshTokenExpiry = 24

    userID := uint(777)
    tenantID := uint(888)

    // Act
    token, err := utils.GenerateRefreshToken(userID, tenantID)
    if err != nil {
        t.Fatalf("GenerateRefreshToken error: %v", err)
    }
    if token == "" {
        t.Fatalf("GenerateRefreshToken returned empty token")
    }

    gotUserID, gotTenantID, err := utils.VerifyRefreshToken(token)
    if err != nil {
        t.Fatalf("VerifyRefreshToken error: %v", err)
    }

    // Assert
    if gotUserID != userID {
        t.Errorf("expected userID %d, got %d", userID, gotUserID)
    }
    if gotTenantID != tenantID {
        t.Errorf("expected tenantID %d, got %d", tenantID, gotTenantID)
    }
}

func TestPasswordHashAndCheck(t *testing.T) {
    password := "s3cr3t-pass"
    hash, err := utils.HashPassword(password)
    if err != nil {
        t.Fatalf("HashPassword error: %v", err)
    }
    if hash == "" {
        t.Fatalf("HashPassword returned empty hash")
    }
    if !utils.CheckPasswordHash(password, hash) {
        t.Errorf("CheckPasswordHash returned false for correct password")
    }
    if utils.CheckPasswordHash("wrong-pass", hash) {
        t.Errorf("CheckPasswordHash returned true for incorrect password")
    }
}