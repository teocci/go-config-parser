// Package parser
// Created by Teocci.
// Author: teocci@yandex.com on 2021-Aug-16
package parser

import (
	"container/list"
	"fmt"
	"github.com/teocci/go-config-parser/core/throw"
	"regexp"
	"strings"
	"sync"
)

// Configuration represents a configuration file with its sections and options.
type Configuration struct {
	filePath        string                // configuration file
	sections        map[string]*list.List // fully qualified section name as key
	orderedSections []string              // track the order of section names as they are parsed
	mutex           sync.RWMutex
}

// NewConfiguration returns a new Configuration instance with an empty file path.
func NewConfiguration() *Configuration {
	return NewConfigurationWithFP(CharEmpty)
}

// NewConfigurationWithFP creates a new Configuration instance.
func NewConfigurationWithFP(filePath string) *Configuration {
	return &Configuration{
		filePath: filePath,
		sections: make(map[string]*list.List),
	}
}

// NewSection creates and adds a new Section with the specified name.
func (c *Configuration) NewSection(fqn string) *Section {
	return c.addSection(fqn)
}

// FilePath returns the configuration file path.
func (c *Configuration) FilePath() string {
	return c.filePath
}

// SetFilePath sets the Configuration file path.
func (c *Configuration) SetFilePath(filePath string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.filePath = filePath
}

// StringValue returns the string value for the specified section and option.
func (c *Configuration) StringValue(section, option string) (value string, err error) {
	s, err := c.Section(section)
	if err != nil {
		return
	}
	value = s.ValueOf(option)
	return
}

// Delete deletes the specified sections matched by a regex name and returns the deleted sections.
func (c *Configuration) Delete(regex string) (sections []*Section, err error) {
	sections, err = c.Find(regex)
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if err == nil {
		for _, s := range sections {
			delete(c.sections, s.Name())
		}
		// remove also from ordered list
		var matched bool
		for i := len(c.orderedSections) - 1; i >= 0; i-- {
			if matched, err = regexp.MatchString(regex, c.orderedSections[i]); matched {
				c.orderedSections = append(c.orderedSections[:i], c.orderedSections[i+1:]...)
			} else {
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return sections, err
}

// Section returns the first section matching the fully qualified section name.
func (c *Configuration) Section(fqn string) (*Section, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if l, ok := c.sections[fqn]; ok {
		for element := l.Front(); element != nil; element = element.Next() {
			s := element.Value.(*Section)
			return s, nil
		}
	}
	return nil, throw.ErrorUnableFindSelection(fqn)
}

// AllSections returns a slice of all sections available.
func (c *Configuration) AllSections() ([]*Section, error) {
	return c.Sections(CharEmpty)
}

// Sections returns a slice of Sections matching the fully qualified section name.
func (c *Configuration) Sections(fqn string) ([]*Section, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	var sections []*Section

	f := func(lst *list.List) {
		for e := lst.Front(); e != nil; e = e.Next() {
			s := e.Value.(*Section)
			sections = append(sections, s)
		}
	}

	if fqn == CharEmpty {
		// Get all sections.
		for _, fqn := range c.orderedSections {
			if lst, ok := c.sections[fqn]; ok {
				f(lst)
			}
		}
	} else {
		if lst, ok := c.sections[fqn]; ok {
			f(lst)
		} else {
			return nil, throw.ErrorUnableFindSelection(fqn)
		}
	}

	return sections, nil
}

// Find returns a slice of Sections matching the regexp against the section name.
func (c *Configuration) Find(regex string) ([]*Section, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	var sections []*Section
	for key, lst := range c.sections {
		if matched, err := regexp.MatchString(regex, key); matched {
			for e := lst.Front(); e != nil; e = e.Next() {
				s := e.Value.(*Section)
				sections = append(sections, s)
			}
		} else {
			if err != nil {
				return nil, err
			}
		}
	}
	return sections, nil
}

// PrintSection prints a text representation of all sections matching the fully qualified section name.
func (c *Configuration) PrintSection(fqn string) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	sections, err := c.Sections(fqn)
	if err == nil {
		for _, section := range sections {
			fmt.Print(section)
		}
	} else {
		fmt.Printf("Unable to find section %v\n", err)
	}
}

// String returns the text representation of a parsed configuration file.
func (c *Configuration) String() string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	var parts []string
	for _, fqn := range c.orderedSections {
		sections, _ := c.Sections(fqn)
		for _, section := range sections {
			parts = append(parts, section.String())
		}
	}
	return strings.Join(parts, CharEmpty)
}

func (c *Configuration) AddSection(fqn string) *Section {
	return c.addSection(fqn)
}

func (c *Configuration) MutexLock() {
	c.mutex.Lock()
}

func (c *Configuration) MutexUnlock() {
	c.mutex.Unlock()
}

func (c *Configuration) addSection(fqn string) *Section {
	section := NewSection(fqn)

	var sectionList *list.List
	if sectionList = c.sections[fqn]; sectionList == nil {
		sectionList = list.New()
		c.sections[fqn] = sectionList
		c.orderedSections = append(c.orderedSections, fqn)
	}

	sectionList.PushBack(section)

	return section
}
