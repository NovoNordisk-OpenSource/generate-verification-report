@URS @reverse_string_feat
Feature: Reverse String
  As a user
  I want to reverse a string via a HTTP GET request
  So that I can read it backwards

  @PV @reverse_string @fixture.client
  Scenario: Reverse a string using a GET request
    Given I have the string "Hello"
    When I send a GET request to reverse it
    Then I should receive "olleH" as the response
