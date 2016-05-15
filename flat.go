package flat

import (
	"bufio"
	"errors"
	"github.com/eaciit/toolkit"
	"io"
	"os"
	"strings"
)

type Flat struct {
	Path       string
	Delimeter  rune
	Ltrim      bool
	Rtrim      bool
	StringMark []string
	Config     toolkit.M

	metadatas       toolkit.M
	file            *os.File
	br              *bufio.Reader
	bw              *bufio.Writer
    lowerName bool
	isRead, isWrite bool
}

type MetaData struct {
	Name  string
	Type  string
	Valid bool
}

func New(path string, r bool, w bool) *Flat {
	f := new(Flat)
	f.Path = path
	f.isRead = w
	f.isWrite = w
	return f
}

func err(pre, txt string) error {
	return errors.New("[Flat." + pre + "] " + txt)
}

func (f *Flat) Open() error {
	var e error
	var file *os.File
	if f.Path == "" {
		return err("Open", "Path is empty")
	}
	file, e = os.Open(f.Path)
	if e != nil {
		return err("Open", e.Error())
	}
	f.file = file

	f.metadatas = toolkit.M{}
	if f.Config.Get("useheader", true).(bool) {
		mheader := toolkit.M{}
		e = f.Read(&mheader)
		if e != nil {
			return err("Open.GetHeader", e.Error())
		}
        
        f.lowerName = f.Config.Get("lowercasename",true).(bool)
		for i, v := range mheader {
            if f.lowerName {
                f.metadatas.Set(i, &MetaData{
                    Name:  strings.ToLower(v.(string)),
                    Type:  "string",
                    Valid: true,
                })
            } else {
                f.metadatas.Set(i, &MetaData{
                    Name:  v.(string),
                    Type:  "string",
                    Valid: true,
                })
            }
		}
	}

	return nil
}

func (f *Flat) Reset() error {
	return nil
}

type MoveFromEnum int

var (
	MoveFromFirst    MoveFromEnum = 0
	MoveFromRelative              = 1
)

func (f *Flat) Move(step int, from MoveFromEnum) error {
	return nil
}

func (f *Flat) MoveRead(obj interface{}, step int, from MoveFromEnum) error {
	if e := f.Move(step, from); e != nil {
		return e
	}
	if e := f.Read(obj); e != nil {
		return e
	}
	return nil
}

func (f *Flat) ReadString()(txt string, e error){
	if f.br == nil {
		f.br = bufio.NewReader(f.file)
	}
	//f.br.ReadLine()
	blines, _, e := f.br.ReadLine()
	if e != nil {
		if e == io.EOF {
			return
		}
		e = err("Read", e.Error())
        return 
	}

	txt = string(blines)
	return
}

func (f *Flat) Read(obj interface{}) error {
	txt, e := f.ReadString()
	if e != nil {
		if e == io.EOF {
			return e
		}
		return err("Read", e.Error())
	}

	m := f.SplitText(txt)
	toolkit.Serde(m, obj, "json")
	return nil
}

func (f *Flat) ReadM() (m toolkit.M, e error) {
	var txt string
    txt, e = f.ReadString()
	if e != nil {
		if e == io.EOF {
			return
		}
		e = err("Read", e.Error())
        return
	}

	m = *f.SplitText(txt)
	return
}

func (f *Flat) SplitText(txt string) *toolkit.M {
	tempStr := ""
	colIndex := 0
	m := &toolkit.M{}
	for _, c := range txt {
		s := string(c)
		if c != f.Delimeter {
			tempStr += s
		} else {
			colIndexStr := toolkit.Sprintf("%d", colIndex)
			metadata := f.metadatas.Get(colIndexStr, new(MetaData)).(*MetaData)
			if metadata.Valid {
				m.Set(metadata.Name, tempStr)
			} else {
				m.Set(toolkit.ToString(colIndex), tempStr)
			}
			tempStr = ""
			colIndex++
		}
	}
	return m
}

func (f *Flat) Write(obj interface{}) error {
	return nil
}

func (f *Flat) Close() error {
	if f.file != nil {
		f.file.Close()
	}
	return nil
}
