package canconv

import (
	"fmt"
	"os"

	"github.com/FerroO2000/canconv/sym"
)

const dbcDefNode = "Vector_XXX"

const dbcHeaders = `
NS_ :
	NS_DESC_
	CM_
	BA_DEF_
	BA_
	VAL_
	CAT_DEF_
	CAT_
	FILTER
	BA_DEF_DEF_
	EV_DATA_
	ENVVAR_DATA_
	SGTYPE_
	SGTYPE_VAL_
	BA_DEF_SGTYPE_
	BA_SGTYPE_
	SIG_TYPE_REF_
	VAL_TABLE_
	SIG_GROUP_
	SIG_VALTYPE_
	SIGTYPE_VALTYPE_
	BO_TX_BU_
	BA_DEF_REL_
	BA_REL_
	BA_DEF_DEF_REL_
	BU_SG_REL_
	BU_EV_REL_
	BU_BO_REL_ 
`

type DBCGenerator struct{}

func NewDBCGenerator() *DBCGenerator {
	return &DBCGenerator{}
}

func (g *DBCGenerator) Generate(model *Model, file *os.File) {
	f := newFile(file)

	f.print("VERSION", FormatString(model.Version))

	f.print(dbcHeaders)

	g.genNodes(f, model.Nodes)
	f.print()

	for msgName, msg := range model.Messages {
		g.genMessage(f, msgName, &msg)
	}
	f.print()

	g.genComments(f, model)
	f.print()

	g.genBitmaps(f, model)
}

func (g *DBCGenerator) genNodes(f *file, nodes map[string]Node) {
	nodeNames := []string{}
	for nodeName := range nodes {
		nodeNames = append(nodeNames, nodeName)
	}

	str := []string{sym.DBCNode}
	str = append(str, nodeNames...)
	f.print(str...)
}

func (g *DBCGenerator) genMessage(f *file, msgName string, msg *Message) {
	id := fmt.Sprintf("%d", msg.ID)
	length := fmt.Sprintf("%d", msg.Length)
	sender := msg.Sender
	if sender == "" {
		sender = dbcDefNode
	}
	f.print(sym.DBCMessage, id, msgName+":", length, sender)

	for sigName, sig := range msg.Signals {
		sig.Validate()
		g.genSignal(f, sigName, &sig)
	}
}

func (g *DBCGenerator) genSignal(f *file, sigName string, sig *Signal) {
	byteOrder := 0
	if sig.BigEndian {
		byteOrder = 1
	}
	valueType := "+"
	if sig.Signed {
		valueType = "-"
	}
	byteDef := fmt.Sprintf("%d|%d@%d%s", sig.StartBit, sig.Size, byteOrder, valueType)
	multiplier := fmt.Sprintf("(%s,%s)", FormatFloat(sig.Scale), FormatFloat(sig.Offset))
	valueRange := fmt.Sprintf("[%s|%s]", FormatFloat(sig.Min), FormatFloat(sig.Max))
	unit := fmt.Sprintf(`"%s"`, sig.Unit)

	recivers := ""
	if len(sig.Recivers) == 0 {
		recivers = dbcDefNode
	} else {
		for i, r := range sig.Recivers {
			if i == 0 {
				recivers += r
				continue
			}
			recivers += "," + r
		}
	}

	f.print("", sym.DBCSignal, sigName, ":", byteDef, multiplier, valueRange, unit, recivers)
}

func (g *DBCGenerator) genComments(f *file, m *Model) {
	for nodeName, node := range m.Nodes {
		if node.HasDescription() {
			f.print(sym.DBCComment, sym.DBCNode, nodeName, FormatString(node.Description), ";")
		}
	}
	f.print()

	for _, msg := range m.Messages {
		if msg.HasDescription() {
			f.print(sym.DBCComment, sym.DBCMessage, msg.FormatID(), FormatString(msg.Description), ";")
		}

		for sigName, sig := range msg.Signals {
			if sig.HasDescription() {
				f.print(sym.DBCComment, sym.DBCSignal, msg.FormatID(), sigName, FormatString(sig.Description), ";")
			}
		}
	}
}

func (g *DBCGenerator) genBitmaps(f *file, m *Model) {
	for _, msg := range m.Messages {
		for sigName, sig := range msg.Signals {
			if sig.IsBitmap() {
				bitmap := ""
				first := true
				for name, val := range sig.Bitmap {
					if first {
						bitmap += FormatUint(val) + " " + FormatString(name)
						first = false
						continue
					}
					bitmap += " " + FormatUint(val) + " " + FormatString(name)
				}
				f.print(sym.DBCValue, msg.FormatID(), sigName, bitmap, ";")
			}
		}
	}
}
