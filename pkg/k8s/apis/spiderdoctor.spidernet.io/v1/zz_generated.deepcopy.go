//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Code generated by controller-gen. DO NOT EDIT.

package v1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Netdoctor) DeepCopyInto(out *Netdoctor) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Netdoctor.
func (in *Netdoctor) DeepCopy() *Netdoctor {
	if in == nil {
		return nil
	}
	out := new(Netdoctor)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Netdoctor) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NetdoctorList) DeepCopyInto(out *NetdoctorList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Netdoctor, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NetdoctorList.
func (in *NetdoctorList) DeepCopy() *NetdoctorList {
	if in == nil {
		return nil
	}
	out := new(NetdoctorList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NetdoctorList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NetdoctorSpec) DeepCopyInto(out *NetdoctorSpec) {
	*out = *in
	if in.Schedule != nil {
		in, out := &in.Schedule, &out.Schedule
		*out = new(SchedulePlan)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NetdoctorSpec.
func (in *NetdoctorSpec) DeepCopy() *NetdoctorSpec {
	if in == nil {
		return nil
	}
	out := new(NetdoctorSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NetdoctorStatus) DeepCopyInto(out *NetdoctorStatus) {
	*out = *in
	if in.DoneRound != nil {
		in, out := &in.DoneRound, &out.DoneRound
		*out = new(int64)
		**out = **in
	}
	if in.LastRoundTimeStamp != nil {
		in, out := &in.LastRoundTimeStamp, &out.LastRoundTimeStamp
		*out = (*in).DeepCopy()
	}
	if in.NextRoundTimeStamp != nil {
		in, out := &in.NextRoundTimeStamp, &out.NextRoundTimeStamp
		*out = (*in).DeepCopy()
	}
	if in.LastRoundStatus != nil {
		in, out := &in.LastRoundStatus, &out.LastRoundStatus
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NetdoctorStatus.
func (in *NetdoctorStatus) DeepCopy() *NetdoctorStatus {
	if in == nil {
		return nil
	}
	out := new(NetdoctorStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SchedulePlan) DeepCopyInto(out *SchedulePlan) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SchedulePlan.
func (in *SchedulePlan) DeepCopy() *SchedulePlan {
	if in == nil {
		return nil
	}
	out := new(SchedulePlan)
	in.DeepCopyInto(out)
	return out
}
