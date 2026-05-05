package auditv1

import proto "github.com/golang/protobuf/proto"

const _ = proto.ProtoPackageIsVersion4

type AnalyzeRequest struct {
	Config   []byte `protobuf:"bytes,1,opt,name=config,proto3" json:"config,omitempty"`
	Filename string `protobuf:"bytes,2,opt,name=filename,proto3" json:"filename,omitempty"`
}

func (m *AnalyzeRequest) Reset() {
	*m = AnalyzeRequest{}
}

func (m *AnalyzeRequest) String() string {
	return proto.CompactTextString(m)
}

func (*AnalyzeRequest) ProtoMessage() {}

func (m *AnalyzeRequest) GetConfig() []byte {
	if m != nil {
		return m.Config
	}
	return nil
}

func (m *AnalyzeRequest) GetFilename() string {
	if m != nil {
		return m.Filename
	}
	return ""
}

type AnalyzeResponse struct {
	Findings []*Finding `protobuf:"bytes,1,rep,name=findings,proto3" json:"findings,omitempty"`
}

func (m *AnalyzeResponse) Reset() {
	*m = AnalyzeResponse{}
}

func (m *AnalyzeResponse) String() string {
	return proto.CompactTextString(m)
}

func (*AnalyzeResponse) ProtoMessage() {}

func (m *AnalyzeResponse) GetFindings() []*Finding {
	if m != nil {
		return m.Findings
	}
	return nil
}

type Finding struct {
	RuleId         string `protobuf:"bytes,1,opt,name=rule_id,json=ruleId,proto3" json:"rule_id,omitempty"`
	Severity       string `protobuf:"bytes,2,opt,name=severity,proto3" json:"severity,omitempty"`
	Message        string `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
	Recommendation string `protobuf:"bytes,4,opt,name=recommendation,proto3" json:"recommendation,omitempty"`
	Path           string `protobuf:"bytes,5,opt,name=path,proto3" json:"path,omitempty"`
	Source         string `protobuf:"bytes,6,opt,name=source,proto3" json:"source,omitempty"`
}

func (m *Finding) Reset() {
	*m = Finding{}
}

func (m *Finding) String() string {
	return proto.CompactTextString(m)
}

func (*Finding) ProtoMessage() {}

func (m *Finding) GetRuleId() string {
	if m != nil {
		return m.RuleId
	}
	return ""
}

func (m *Finding) GetSeverity() string {
	if m != nil {
		return m.Severity
	}
	return ""
}

func (m *Finding) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *Finding) GetRecommendation() string {
	if m != nil {
		return m.Recommendation
	}
	return ""
}

func (m *Finding) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

func (m *Finding) GetSource() string {
	if m != nil {
		return m.Source
	}
	return ""
}
