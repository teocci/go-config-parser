// Package parser
// Created by Teocci.
// Author: teocci@yandex.com on 2021-Aug-15
package parser

import (
	"bufio"
	"os"
	"path"
	"strings"
)

// CharEqual is the delimiter to be used between section key and values.
const (
	CharEmpty        = ""
	CharSpace        = " "
	CharEqual        = "="
	CharColon        = ":"
	CharLSBracket    = "["
	CharRSBracket    = "]"
	CharEOF          = "\n"
	SectionEmpty     = CharLSBracket + CharRSBracket
	BackupExt        = ".bak"
)

type Delimiter int

const (
	EqualDelimiter Delimiter = iota
	ColonDelimiter
)

const (
	EqualDelimiterString = CharSpace + CharEqual + CharSpace
	ColonDelimiterString = CharSpace + CharColon + CharSpace
)

var delimiter = EqualDelimiter

func DelimiterString(d Delimiter) string {
	return delimiterStrings()[d]
}

func delimiterStrings() []string {
	return []string{EqualDelimiterString, ColonDelimiterString}
}

func SetDelimiter(d Delimiter) {
	delimiter = d
}

// Read parses a specified configuration file and returns a Configuration instance.
func Read(filePath string) (*Configuration, error) {
	filePath = path.Clean(filePath)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := NewConfigurationWithFP(filePath)
	activeSection := config.AddSection(globalSection)

	scanner := bufio.NewScanner(bufio.NewReader(file))
	for scanner.Scan() {
		// TODO: maybe trim spaces here
		line := scanner.Text()
		if len(line) < 0 {
			continue
		}

		if IsSection(line) {
			fqn := strings.Trim(line, CharSpace+SectionEmpty)
			activeSection = config.AddSection(fqn)
			continue
		}
		// save options and comments
		AddOption(activeSection, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return config, nil
}

// Save the Configuration to file. Creates a backup (.bak) if file already exists.
func Save(c *Configuration, filePath string) (err error) {
	c.MutexLock()

	//fine if the file does not exist
	if err := os.Rename(filePath, filePath + BackupExt); err != nil && !os.IsNotExist(err) {
		return err
	}

	f, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer func() {
		err = f.Close()
	}()

	w := bufio.NewWriter(f)
	defer func() {
		err = w.Flush()
	}()
	c.MutexUnlock()

	s, err := c.AllSections()
	if err != nil {
		return err
	}

	c.MutexLock()
	defer c.MutexUnlock()

	for _, v := range s {
		_, err := w.WriteString(v.String())
		if err != nil {
			return err
		}
	}

	return err
}
