//go:build !js

package main

func showPadsFromURL() bool { return false }

func pollPhoneKeys() phoneKeys { return phoneKeys{} }

func editFromURL() bool { return false }
