Feature: Member login via email and password
  In order to access protected resources
  As a registered member
  I want to log in and receive access and refresh tokens

  Background:
    Given the auth API server is running
    And the following members exist:
      | email            | password    |
      | user@example.com | password123 |

  Scenario: Successful login with valid email and password
    When I POST "/auth/login" with headers:
      | User-Agent      | test-agent   |
      | X-Forwarded-For | 127.0.0.1    |
    And the JSON body:
      """
      {
        "email": "user@example.com",
        "password": "password123"
      }
      """
    Then the response status code should be 200
    And the response JSON should contain a non-empty "access_token"
    And the response JSON should contain a non-empty "refresh_token"
    And the response JSON field "refresh_max_age_sec" should equal 3600
    And the response JSON field "refresh_endpoint" should equal "/auth/refresh"
    And a refresh token session should be persisted for member "user@example.com" from user agent "test-agent" and IP "127.0.0.1"

  Scenario: Login fails with wrong password
    When I POST "/auth/login" with headers:
      | User-Agent      | test-agent   |
      | X-Forwarded-For | 127.0.0.1    |
    And the JSON body:
      """
      {
        "email": "user@example.com",
        "password": "wrong-password"
      }
      """
    Then the response status code should be 401
    And the response JSON field "error_code" should equal "invalid_credentials"
    And the response JSON field "message" should not be empty
    And no new refresh token session should be persisted for member "user@example.com"

  Scenario: Login fails for unknown email
    When I POST "/auth/login" with headers:
      | User-Agent      | test-agent   |
      | X-Forwarded-For | 127.0.0.1    |
    And the JSON body:
      """
      {
        "email": "unknown@example.com",
        "password": "password123"
      }
      """
    Then the response status code should be 401
    And the response JSON field "error_code" should equal "invalid_credentials"
    And the response JSON field "message" should not be empty
    And no refresh token session should be persisted for member "unknown@example.com"

  Scenario: Login fails when token service is unavailable
    Given the access token service is temporarily unavailable
    When I POST "/auth/login" with headers:
      | User-Agent      | test-agent   |
      | X-Forwarded-For | 127.0.0.1    |
    And the JSON body:
      """
      {
        "email": "user@example.com",
        "password": "password123"
      }
      """
    Then the response status code should be 500
    And the response JSON field "error_code" should equal "internal_error"
    And no refresh token session should be persisted for member "user@example.com"

  Scenario: Login fails with invalid request body
    When I POST "/auth/login" with headers:
      | User-Agent      | test-agent   |
      | X-Forwarded-For | 127.0.0.1    |
    And the raw body:
      """
      { "email": 123, "password": true }
      """
    Then the response status code should be 400
    And the response JSON field "error_code" should equal "invalid_request"
    And the response JSON field "message" should not be empty
    And no refresh token session should be persisted
