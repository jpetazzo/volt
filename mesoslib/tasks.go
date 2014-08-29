package mesoslib

import (
	"fmt"
	"strings"

	"code.google.com/p/goprotobuf/proto"
	"github.com/Sirupsen/logrus"
	"github.com/VoltFramework/volt/mesosproto"
)

func createTaskInfo(offer *mesosproto.Offer, resources []*mesosproto.Resource, args []string, ID, image string) *mesosproto.TaskInfo {
	taskInfo := mesosproto.TaskInfo{
		Name: proto.String(fmt.Sprintf("volt-task-%s", ID)),
		TaskId: &mesosproto.TaskID{
			Value: &ID,
		},
		SlaveId:   offer.SlaveId,
		Resources: resources,
		Command:   &mesosproto.CommandInfo{},
	}

	// Set value only if provided
	if args[0] != "" {
		taskInfo.Command.Value = &args[0]
	}

	// Set args only if they exist
	fmt.Println(len(args))
	if len(args) > 1 {
		taskInfo.Command.Arguments = args[1:]
	}

	// Set the docker image if specified
	if image != "" {
		taskInfo.Container = &mesosproto.ContainerInfo{
			Type: mesosproto.ContainerInfo_DOCKER.Enum(),
			Docker: &mesosproto.ContainerInfo_DockerInfo{
				Image: &image,
			},
		}
		taskInfo.Command.Shell = proto.Bool(false)
	}
	return &taskInfo
}

func (m *MesosLib) LaunchTask(offer *mesosproto.Offer, resources []*mesosproto.Resource, command, ID, image string) error {
	m.Log.WithFields(logrus.Fields{"ID": ID, "command": command, "offerId": offer.Id, "dockerImage": image}).Info("Launching task...")

	taskInfo := createTaskInfo(offer, resources, strings.Split(command, " "), ID, image)

	return m.send(&mesosproto.LaunchTasksMessage{
		FrameworkId: m.frameworkInfo.Id,
		Tasks:       []*mesosproto.TaskInfo{taskInfo},
		OfferIds: []*mesosproto.OfferID{
			offer.Id,
		},
		Filters: &mesosproto.Filters{},
	}, "mesos.internal.LaunchTasksMessage")
}

func (m *MesosLib) KillTask(ID string) error {
	m.Log.WithFields(logrus.Fields{"ID": ID}).Info("Killing task...")

	return m.send(&mesosproto.KillTaskMessage{
		FrameworkId: m.frameworkInfo.Id,
		TaskId: &mesosproto.TaskID{
			Value: &ID,
		},
	}, "mesos.internal.KillTaskMessage")
}
