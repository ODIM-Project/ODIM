package rpc

import (
	"context"
	"errors"

	accountproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/account"
	aggregatorproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/aggregator"
	chassisproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/chassis"
	"github.com/ODIM-Project/ODIM/lib-utilities/proto/events"
	eventsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/events"
	fabricsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/fabrics"
	managersproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/managers"
	roleproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/role"
	sessionproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/session"
	systemsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/systems"
	taskproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/task"
	teleproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/telemetry"
	updateproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/update"
	"google.golang.org/grpc"
)

type fakeStruct struct{}
type fakeStruct2 struct{}

// -------------------------ACCOUNT-----------------------------------------------
func (fakeStruct) GetAccountServices(ctx context.Context, in *accountproto.AccountRequest, opts ...grpc.CallOption) (*accountproto.AccountResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) Create(ctx context.Context, in *accountproto.CreateAccountRequest, opts ...grpc.CallOption) (*accountproto.AccountResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) GetAllAccounts(ctx context.Context, in *accountproto.AccountRequest, opts ...grpc.CallOption) (*accountproto.AccountResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) GetAccount(ctx context.Context, in *accountproto.GetAccountRequest, opts ...grpc.CallOption) (*accountproto.AccountResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) Update(ctx context.Context, in *accountproto.UpdateAccountRequest, opts ...grpc.CallOption) (*accountproto.AccountResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) Delete(ctx context.Context, in *accountproto.DeleteAccountRequest, opts ...grpc.CallOption) (*accountproto.AccountResponse, error) {
	return nil, errors.New("fakeError")
}

//------------------------------------AGGREGATOR-------------------------------------------------

func (fakeStruct) Reset(ctx context.Context, in *aggregatorproto.AggregatorRequest, opts ...grpc.CallOption) (*aggregatorproto.AggregatorResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) GetAggregationService(ctx context.Context, in *aggregatorproto.AggregatorRequest, opts ...grpc.CallOption) (*aggregatorproto.AggregatorResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) SetDefaultBootOrder(ctx context.Context, in *aggregatorproto.AggregatorRequest, opts ...grpc.CallOption) (*aggregatorproto.AggregatorResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) RediscoverSystemInventory(ctx context.Context, in *aggregatorproto.RediscoverSystemInventoryRequest, opts ...grpc.CallOption) (*aggregatorproto.RediscoverSystemInventoryResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) UpdateSystemState(ctx context.Context, in *aggregatorproto.UpdateSystemStateRequest, opts ...grpc.CallOption) (*aggregatorproto.UpdateSystemStateResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) AddAggregationSource(ctx context.Context, in *aggregatorproto.AggregatorRequest, opts ...grpc.CallOption) (*aggregatorproto.AggregatorResponse, error) {

	return nil, errors.New("fakeError")
}

func (fakeStruct) GetAllAggregationSource(ctx context.Context, in *aggregatorproto.AggregatorRequest, opts ...grpc.CallOption) (*aggregatorproto.AggregatorResponse, error) {

	return nil, errors.New("fakeError")
}

func (fakeStruct) GetAggregationSource(ctx context.Context, in *aggregatorproto.AggregatorRequest, opts ...grpc.CallOption) (*aggregatorproto.AggregatorResponse, error) {

	return nil, errors.New("fakeError")
}

func (fakeStruct) UpdateAggregationSource(ctx context.Context, in *aggregatorproto.AggregatorRequest, opts ...grpc.CallOption) (*aggregatorproto.AggregatorResponse, error) {

	return nil, errors.New("fakeError")
}

func (fakeStruct) DeleteAggregationSource(ctx context.Context, in *aggregatorproto.AggregatorRequest, opts ...grpc.CallOption) (*aggregatorproto.AggregatorResponse, error) {

	return nil, errors.New("fakeError")
}

func (fakeStruct) CreateAggregate(ctx context.Context, in *aggregatorproto.AggregatorRequest, opts ...grpc.CallOption) (*aggregatorproto.AggregatorResponse, error) {

	return nil, errors.New("fakeError")
}

func (fakeStruct) GetAllAggregates(ctx context.Context, in *aggregatorproto.AggregatorRequest, opts ...grpc.CallOption) (*aggregatorproto.AggregatorResponse, error) {

	return nil, errors.New("fakeError")
}

func (fakeStruct) GetAggregate(ctx context.Context, in *aggregatorproto.AggregatorRequest, opts ...grpc.CallOption) (*aggregatorproto.AggregatorResponse, error) {

	return nil, errors.New("fakeError")
}

func (fakeStruct) DeleteAggregate(ctx context.Context, in *aggregatorproto.AggregatorRequest, opts ...grpc.CallOption) (*aggregatorproto.AggregatorResponse, error) {

	return nil, errors.New("fakeError")
}

func (fakeStruct) AddElementsToAggregate(ctx context.Context, in *aggregatorproto.AggregatorRequest, opts ...grpc.CallOption) (*aggregatorproto.AggregatorResponse, error) {

	return nil, errors.New("fakeError")
}

func (fakeStruct) RemoveElementsFromAggregate(ctx context.Context, in *aggregatorproto.AggregatorRequest, opts ...grpc.CallOption) (*aggregatorproto.AggregatorResponse, error) {

	return nil, errors.New("fakeError")
}

func (fakeStruct) ResetElementsOfAggregate(ctx context.Context, in *aggregatorproto.AggregatorRequest, opts ...grpc.CallOption) (*aggregatorproto.AggregatorResponse, error) {

	return nil, errors.New("fakeError")
}

func (fakeStruct) SetDefaultBootOrderElementsOfAggregate(ctx context.Context, in *aggregatorproto.AggregatorRequest, opts ...grpc.CallOption) (*aggregatorproto.AggregatorResponse, error) {

	return nil, errors.New("fakeError")
}

func (fakeStruct) GetAllConnectionMethods(ctx context.Context, in *aggregatorproto.AggregatorRequest, opts ...grpc.CallOption) (*aggregatorproto.AggregatorResponse, error) {

	return nil, errors.New("fakeError")
}

func (fakeStruct) GetConnectionMethod(ctx context.Context, in *aggregatorproto.AggregatorRequest, opts ...grpc.CallOption) (*aggregatorproto.AggregatorResponse, error) {

	return nil, errors.New("fakeError")
}

func (fakeStruct) SendStartUpData(ctx context.Context, in *aggregatorproto.SendStartUpDataRequest, opts ...grpc.CallOption) (*aggregatorproto.SendStartUpDataResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) GetResetActionInfoService(ctx context.Context, in *aggregatorproto.AggregatorRequest, opts ...grpc.CallOption) (*aggregatorproto.AggregatorResponse, error) {

	return nil, errors.New("fakeError")
}

func (fakeStruct) GetSetDefaultBootOrderActionInfo(ctx context.Context, in *aggregatorproto.AggregatorRequest, opts ...grpc.CallOption) (*aggregatorproto.AggregatorResponse, error) {

	return nil, errors.New("fakeError")
}

func (fakeStruct) IsAggregateHaveSubscription(ctx context.Context, in *events.EventUpdateRequest, opts ...grpc.CallOption) (*events.SubscribeEMBResponse, error) {

	return nil, errors.New("fakeError")
}

func (fakeStruct) RemoveEventSubscriptionsRPC(ctx context.Context, in *events.EventUpdateRequest, opts ...grpc.CallOption) (*events.SubscribeEMBResponse, error) {

	return nil, errors.New("fakeError")
}
func (fakeStruct) UpdateEventSubscriptionsRPC(ctx context.Context, in *events.EventUpdateRequest, opts ...grpc.CallOption) (*events.SubscribeEMBResponse, error) {

	return nil, errors.New("fakeError")
}

//--------------------------------CHASSIS--------------------------------

func (fakeStruct) GetChassisCollection(ctx context.Context, in *chassisproto.GetChassisRequest, opts ...grpc.CallOption) (*chassisproto.GetChassisResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) GetChassisResource(ctx context.Context, in *chassisproto.GetChassisRequest, opts ...grpc.CallOption) (*chassisproto.GetChassisResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) GetChassisInfo(ctx context.Context, in *chassisproto.GetChassisRequest, opts ...grpc.CallOption) (*chassisproto.GetChassisResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) CreateChassis(ctx context.Context, in *chassisproto.CreateChassisRequest, opts ...grpc.CallOption) (*chassisproto.GetChassisResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) DeleteChassis(ctx context.Context, in *chassisproto.DeleteChassisRequest, opts ...grpc.CallOption) (*chassisproto.GetChassisResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) UpdateChassis(ctx context.Context, in *chassisproto.UpdateChassisRequest, opts ...grpc.CallOption) (*chassisproto.GetChassisResponse, error) {
	return nil, errors.New("fakeError")
}

//-------------------------------------EVENTS------------------------------------

func (fakeStruct) GetEventService(ctx context.Context, in *eventsproto.EventSubRequest, opts ...grpc.CallOption) (*eventsproto.EventSubResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) SubmitTestEvent(ctx context.Context, in *eventsproto.EventSubRequest, opts ...grpc.CallOption) (*eventsproto.EventSubResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) CreateEventSubscription(ctx context.Context, in *eventsproto.EventSubRequest, opts ...grpc.CallOption) (*eventsproto.EventSubResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) GetEventSubscription(ctx context.Context, in *eventsproto.EventRequest, opts ...grpc.CallOption) (*eventsproto.EventSubResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) DeleteEventSubscription(ctx context.Context, in *eventsproto.EventRequest, opts ...grpc.CallOption) (*eventsproto.EventSubResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) CreateDefaultEventSubscription(ctx context.Context, in *eventsproto.DefaultEventSubRequest, opts ...grpc.CallOption) (*eventsproto.DefaultEventSubResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) GetEventSubscriptionsCollection(ctx context.Context, in *eventsproto.EventRequest, opts ...grpc.CallOption) (*eventsproto.EventSubResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) SubsribeEMB(ctx context.Context, in *eventsproto.SubscribeEMBRequest, opts ...grpc.CallOption) (*eventsproto.SubscribeEMBResponse, error) {
	return nil, errors.New("fakeError")
}

//--------------------------------------FABRICS--------------------------------------

func (fakeStruct) GetFabricResource(ctx context.Context, in *fabricsproto.FabricRequest, opts ...grpc.CallOption) (*fabricsproto.FabricResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) UpdateFabricResource(ctx context.Context, in *fabricsproto.FabricRequest, opts ...grpc.CallOption) (*fabricsproto.FabricResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) AddFabric(ctx context.Context, in *fabricsproto.AddFabricRequest, opts ...grpc.CallOption) (*fabricsproto.FabricResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) DeleteFabricResource(ctx context.Context, in *fabricsproto.FabricRequest, opts ...grpc.CallOption) (*fabricsproto.FabricResponse, error) {
	return nil, errors.New("fakeError")
}
func (fakeStruct) RemoveFabric(ctx context.Context, in *fabricsproto.AddFabricRequest, opts ...grpc.CallOption) (*fabricsproto.FabricResponse, error) {
	return nil, errors.New("fakeError")
}

//-----------------------------------MANAGERS-----------------------------------------------

func (fakeStruct) GetManagersCollection(ctx context.Context, in *managersproto.ManagerRequest, opts ...grpc.CallOption) (*managersproto.ManagerResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) GetManager(ctx context.Context, in *managersproto.ManagerRequest, opts ...grpc.CallOption) (*managersproto.ManagerResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) GetManagersResource(ctx context.Context, in *managersproto.ManagerRequest, opts ...grpc.CallOption) (*managersproto.ManagerResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) VirtualMediaInsert(ctx context.Context, in *managersproto.ManagerRequest, opts ...grpc.CallOption) (*managersproto.ManagerResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) VirtualMediaEject(ctx context.Context, in *managersproto.ManagerRequest, opts ...grpc.CallOption) (*managersproto.ManagerResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) GetRemoteAccountService(ctx context.Context, in *managersproto.ManagerRequest, opts ...grpc.CallOption) (*managersproto.ManagerResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) CreateRemoteAccountService(ctx context.Context, in *managersproto.ManagerRequest, opts ...grpc.CallOption) (*managersproto.ManagerResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) UpdateRemoteAccountService(ctx context.Context, in *managersproto.ManagerRequest, opts ...grpc.CallOption) (*managersproto.ManagerResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) DeleteRemoteAccountService(ctx context.Context, in *managersproto.ManagerRequest, opts ...grpc.CallOption) (*managersproto.ManagerResponse, error) {
	return nil, errors.New("fakeError")
}

//------------------------------------ROLE-------------------------------------------------

func (fakeStruct) CreateRole(ctx context.Context, in *roleproto.RoleRequest, opts ...grpc.CallOption) (*roleproto.RoleResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) GetRole(ctx context.Context, in *roleproto.GetRoleRequest, opts ...grpc.CallOption) (*roleproto.RoleResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) GetAllRoles(ctx context.Context, in *roleproto.GetRoleRequest, opts ...grpc.CallOption) (*roleproto.RoleResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) DeleteRole(ctx context.Context, in *roleproto.DeleteRoleRequest, opts ...grpc.CallOption) (*roleproto.RoleResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) UpdateRole(ctx context.Context, in *roleproto.UpdateRoleRequest, opts ...grpc.CallOption) (*roleproto.RoleResponse, error) {
	return nil, errors.New("fakeError")
}

//-------------------------------------SESSIONS---------------------------------------------

func (fakeStruct) CreateSession(ctx context.Context, in *sessionproto.SessionCreateRequest, opts ...grpc.CallOption) (*sessionproto.SessionCreateResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) DeleteSession(ctx context.Context, in *sessionproto.SessionRequest, opts ...grpc.CallOption) (*sessionproto.SessionResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) GetAllActiveSessions(ctx context.Context, in *sessionproto.SessionRequest, opts ...grpc.CallOption) (*sessionproto.SessionResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) GetSession(ctx context.Context, in *sessionproto.SessionRequest, opts ...grpc.CallOption) (*sessionproto.SessionResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) GetSessionUserName(ctx context.Context, in *sessionproto.SessionRequest, opts ...grpc.CallOption) (*sessionproto.SessionUserName, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) GetSessionService(ctx context.Context, in *sessionproto.SessionRequest, opts ...grpc.CallOption) (*sessionproto.SessionResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) GetSessionUserRoleID(ctx context.Context, in *sessionproto.SessionRequest, opts ...grpc.CallOption) (*sessionproto.SessionUsersRoleID, error) {
	return nil, errors.New("fakeError")
}

//--------------------------------------------SYSTEM-----------------------------------------

func (fakeStruct2) GetSystemsCollection(ctx context.Context, in *systemsproto.GetSystemsRequest, opts ...grpc.CallOption) (*systemsproto.SystemsResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct2) GetSystemResource(ctx context.Context, in *systemsproto.GetSystemsRequest, opts ...grpc.CallOption) (*systemsproto.SystemsResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct2) GetSystems(ctx context.Context, in *systemsproto.GetSystemsRequest, opts ...grpc.CallOption) (*systemsproto.SystemsResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct2) ComputerSystemReset(ctx context.Context, in *systemsproto.ComputerSystemResetRequest, opts ...grpc.CallOption) (*systemsproto.SystemsResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct2) SetDefaultBootOrder(ctx context.Context, in *systemsproto.DefaultBootOrderRequest, opts ...grpc.CallOption) (*systemsproto.SystemsResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct2) ChangeBiosSettings(ctx context.Context, in *systemsproto.BiosSettingsRequest, opts ...grpc.CallOption) (*systemsproto.SystemsResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct2) ChangeBootOrderSettings(ctx context.Context, in *systemsproto.BootOrderSettingsRequest, opts ...grpc.CallOption) (*systemsproto.SystemsResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct2) CreateVolume(ctx context.Context, in *systemsproto.VolumeRequest, opts ...grpc.CallOption) (*systemsproto.SystemsResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct2) DeleteVolume(ctx context.Context, in *systemsproto.VolumeRequest, opts ...grpc.CallOption) (*systemsproto.SystemsResponse, error) {
	return nil, errors.New("fakeError")
}

//-----------------------------------------TASK------------------------------------------

func (fakeStruct) DeleteTask(ctx context.Context, in *taskproto.GetTaskRequest, opts ...grpc.CallOption) (*taskproto.TaskResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) GetTasks(ctx context.Context, in *taskproto.GetTaskRequest, opts ...grpc.CallOption) (*taskproto.TaskResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) GetSubTasks(ctx context.Context, in *taskproto.GetTaskRequest, opts ...grpc.CallOption) (*taskproto.TaskResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) GetSubTask(ctx context.Context, in *taskproto.GetTaskRequest, opts ...grpc.CallOption) (*taskproto.TaskResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) TaskCollection(ctx context.Context, in *taskproto.GetTaskRequest, opts ...grpc.CallOption) (*taskproto.TaskResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) GetTaskService(ctx context.Context, in *taskproto.GetTaskRequest, opts ...grpc.CallOption) (*taskproto.TaskResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) GetTaskMonitor(ctx context.Context, in *taskproto.GetTaskRequest, opts ...grpc.CallOption) (*taskproto.TaskResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) CreateTask(ctx context.Context, in *taskproto.CreateTaskRequest, opts ...grpc.CallOption) (*taskproto.CreateTaskResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) CreateChildTask(ctx context.Context, in *taskproto.CreateTaskRequest, opts ...grpc.CallOption) (*taskproto.CreateTaskResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) UpdateTask(ctx context.Context, in *taskproto.UpdateTaskRequest, opts ...grpc.CallOption) (*taskproto.UpdateTaskResponse, error) {
	return nil, errors.New("fakeError")
}

//------------------------------------------TELEMETRY---------------------------------------

func (fakeStruct) GetTelemetryService(ctx context.Context, in *teleproto.TelemetryRequest, opts ...grpc.CallOption) (*teleproto.TelemetryResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) GetMetricDefinitionCollection(ctx context.Context, in *teleproto.TelemetryRequest, opts ...grpc.CallOption) (*teleproto.TelemetryResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) GetMetricReportDefinitionCollection(ctx context.Context, in *teleproto.TelemetryRequest, opts ...grpc.CallOption) (*teleproto.TelemetryResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) GetMetricReportCollection(ctx context.Context, in *teleproto.TelemetryRequest, opts ...grpc.CallOption) (*teleproto.TelemetryResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) GetTriggerCollection(ctx context.Context, in *teleproto.TelemetryRequest, opts ...grpc.CallOption) (*teleproto.TelemetryResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) GetMetricDefinition(ctx context.Context, in *teleproto.TelemetryRequest, opts ...grpc.CallOption) (*teleproto.TelemetryResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) GetMetricReportDefinition(ctx context.Context, in *teleproto.TelemetryRequest, opts ...grpc.CallOption) (*teleproto.TelemetryResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) GetMetricReport(ctx context.Context, in *teleproto.TelemetryRequest, opts ...grpc.CallOption) (*teleproto.TelemetryResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) GetTrigger(ctx context.Context, in *teleproto.TelemetryRequest, opts ...grpc.CallOption) (*teleproto.TelemetryResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) UpdateTrigger(ctx context.Context, in *teleproto.TelemetryRequest, opts ...grpc.CallOption) (*teleproto.TelemetryResponse, error) {
	return nil, errors.New("fakeError")
}

//--------------------------------------------UPDATE----------------------------------------

func (fakeStruct) GetUpdateService(ctx context.Context, in *updateproto.UpdateRequest, opts ...grpc.CallOption) (*updateproto.UpdateResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) GetFirmwareInventory(ctx context.Context, in *updateproto.UpdateRequest, opts ...grpc.CallOption) (*updateproto.UpdateResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) GetFirmwareInventoryCollection(ctx context.Context, in *updateproto.UpdateRequest, opts ...grpc.CallOption) (*updateproto.UpdateResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) GetSoftwareInventory(ctx context.Context, in *updateproto.UpdateRequest, opts ...grpc.CallOption) (*updateproto.UpdateResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) GetSoftwareInventoryCollection(ctx context.Context, in *updateproto.UpdateRequest, opts ...grpc.CallOption) (*updateproto.UpdateResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) SimepleUpdate(ctx context.Context, in *updateproto.UpdateRequest, opts ...grpc.CallOption) (*updateproto.UpdateResponse, error) {
	return nil, errors.New("fakeError")
}

func (fakeStruct) StartUpdate(ctx context.Context, in *updateproto.UpdateRequest, opts ...grpc.CallOption) (*updateproto.UpdateResponse, error) {
	return nil, errors.New("fakeError")
}
