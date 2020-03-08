// +build !ignore_autogenerated

package v1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

func (in *Mydemo) DeepCopyInto(out *Mydemo) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	return
}

func (in *Mydemo) DeepCopy() *Mydemo {
	if in == nil {
		return nil
	}
	out := new(Mydemo)
	in.DeepCopyInto(out)
	return out
}

func (in *Mydemo) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *MydemoList) DeepCopyInto(out *MydemoList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Mydemo, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

func (in *MydemoList) DeepCopy() *MydemoList {
	if in == nil {
		return nil
	}
	out := new(MydemoList)
	in.DeepCopyInto(out)
	return out
}

func (in *MydemoList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *MydemoSpec) DeepCopyInto(out *MydemoSpec) {
	*out = *in
	return
}

func (in *MydemoSpec) DeepCopy() *MydemoSpec {
	if in == nil {
		return nil
	}
	out := new(MydemoSpec)
	in.DeepCopyInto(out)
	return out
}