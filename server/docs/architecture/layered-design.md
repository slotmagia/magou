# 分层设计详解

## 分层架构概述

本项目采用经典的分层架构设计，从上到下分为以下几层：

1. **路由层(Router)**: 负责 URL 路由映射
2. **控制器层(Controller)**: 处理 HTTP 请求和响应
3. **服务层(Service)**: 定义业务接口抽象
4. **逻辑层(Logic)**: 实现具体的业务逻辑
5. **模型层(Model)**: 定义数据结构和实体
6. **数据层(ORM)**: 数据库操作抽象

## 各层职责详解

### 路由层 (Router)

负责将 HTTP 请求映射到对应的控制器方法，主要包含 URL 路径和 HTTP 方法的映射关系。

```go
// 路由注册示例
func Api(ctx context.Context, group *ghttp.RouterGroup) {
    group.Group(simple.RouterPrefix(ctx, consts.AppApi), func(group *ghttp.RouterGroup) {
        group.Bind(
            api.Role(), // 角色管理路由
            api.User(), // 用户管理路由
        )
        group.Middleware(service.Middleware().ApiAuth)
    })
}
```

### 控制器层 (Controller)

负责接收 HTTP 请求、参数校验、调用服务层接口、返回响应结果。控制器方法与 API 接口一一对应。

```go
// 控制器示例 - 用户登录
func (c *cUser) Login(ctx context.Context, req *user.LoginReq) (res *user.LoginRes, err error) {
    out, err := service.User().Login(ctx, &req.UserLoginInp)
    if err != nil {
        return nil, err
    }
    return &user.LoginRes{
        LoginTokenModel: out,
    }, nil
}
```

### 服务层 (Service)

定义业务接口抽象，是控制器层和逻辑层之间的桥梁。服务层只包含接口定义，不包含具体实现。

```go
// 用户服务接口示例（实际定义）
type IUser interface {
    // Login 用户登录
    Login(ctx context.Context, in *sysin.UserLoginInp) (res *sysout.LoginTokenModel, err error)
    
    // Logout 用户退出
    Logout(ctx context.Context, in *sysin.UserLogoutInp) error
    
    // GetProfile 获取用户资料
    GetProfile(ctx context.Context, userId int64) (res *sysout.UserModel, err error)
    
    // RefreshToken 刷新访问令牌
    RefreshToken(ctx context.Context, refreshToken string) (res *TokenInfo, err error)
    
    // ChangePassword 修改密码
    ChangePassword(ctx context.Context, userId int64, oldPassword, newPassword string) error
    
    // GenerateCaptcha 生成验证码
    GenerateCaptcha(ctx context.Context) (captchaId, captchaImage string, err error)
    
    // VerifyCaptcha 验证验证码
    VerifyCaptcha(ctx context.Context, captchaId, captcha string) error
    
    // GetUserByUsername 根据用户名获取用户
    GetUserByUsername(ctx context.Context, username string) (user *sysout.UserModel, err error)
    
    // ValidateUser 验证用户密码
    ValidateUser(ctx context.Context, username, password string) (user *sysout.UserModel, err error)
}

// 角色服务接口示例（实际定义）
type IRole interface {
    // 角色基础操作
    GetRoleList(ctx context.Context, in *sysin.RoleListInp) (res *sysout.RoleListModel, err error)
    GetRoleDetail(ctx context.Context, in *sysin.RoleDetailInp) (res *sysout.RoleDetailModel, err error)
    CreateRole(ctx context.Context, in *sysin.CreateRoleInp) (res *sysout.RoleModel, err error)
    UpdateRole(ctx context.Context, in *sysin.UpdateRoleInp) (res *sysout.RoleModel, err error)
    DeleteRole(ctx context.Context, in *sysin.DeleteRoleInp) (err error)
    BatchDeleteRole(ctx context.Context, in *sysin.BatchDeleteRoleInp) (err error)
    UpdateRoleStatus(ctx context.Context, in *sysin.UpdateRoleStatusInp) (err error)
    CopyRole(ctx context.Context, in *sysin.CopyRoleInp) (res *sysout.RoleModel, err error)

    // 角色权限管理
    GetRoleMenus(ctx context.Context, in *sysin.RoleMenuInp) (res *sysout.RoleMenuModel, err error)
    UpdateRoleMenus(ctx context.Context, in *sysin.UpdateRoleMenuInp) (err error)
    GetRolePermissions(ctx context.Context, in *sysin.RolePermissionInp) (res *sysout.RolePermissionModel, err error)

    // 角色选项和统计
    GetRoleOptions(ctx context.Context, in *sysin.RoleOptionInp) (res []*sysout.RoleOptionModel, err error)
    GetRoleStats(ctx context.Context) (res *sysout.RoleStatsModel, err error)
    GetDataScopeOptions(ctx context.Context) (res []*sysout.DataScopeModel, err error)

    // 用户角色关联操作
    AssignUserRoles(ctx context.Context, userId int64, roleIds []int64, assignedBy int64) (err error)
    RemoveUserRoles(ctx context.Context, userId int64, roleIds []int64) (err error)
    GetUserRoles(ctx context.Context, userId int64) (res []*sysout.RoleModel, err error)
    SetUserPrimaryRole(ctx context.Context, userId int64, roleId int64) (err error)

    // 权限验证
    CheckUserPermission(ctx context.Context, userId int64, permission string) (bool, error)
    CheckUserRole(ctx context.Context, userId int64, roleCode string) (bool, error)
    GetUserPermissions(ctx context.Context, userId int64) ([]string, error)
    GetUserMenus(ctx context.Context, userId int64) ([]int64, error)
    GetUserDataScope(ctx context.Context, userId int64) (int, error)

    // 批量权限验证
    CheckUsersPermission(ctx context.Context, userIds []int64, permission string) (map[int64]bool, error)
    FilterUsersByPermission(ctx context.Context, userIds []int64, permission string) ([]int64, error)
}

// 服务获取函数
func User() IUser {
    if localUser == nil {
        panic("implement not found for interface IUser, forgot register?")
    }
    return localUser
}

func Role() IRole {
    if localRole == nil {
        panic("implement not found for interface IRole, forgot register?")
    }
    return localRole
}
```

### 逻辑层 (Logic)

包含具体的业务逻辑实现，实现服务层定义的接口。所有业务规则、数据库操作和流程都在此层实现。

```go
// 用户业务逻辑实现示例（实际代码）
type sUser struct{}

func init() {
    service.RegisterUser(NewUser())
}

func NewUser() *sUser {
    return &sUser{}
}

// Login 用户登录实现
func (s *sUser) Login(ctx context.Context, in *sysin.UserLoginInp) (res *sysout.LoginTokenModel, err error) {
    // 验证用户密码
    user, err := s.ValidateUser(ctx, in.Username, in.Password)
    if err != nil {
        return nil, err
    }

    // 检查用户状态
    if user.Status != entity.UserStatusNormal {
        switch user.Status {
        case entity.UserStatusLocked:
            return nil, gerror.New(consts.GetAuthErrorMessage(consts.ErrUserLocked))
        case entity.UserStatusDisabled:
            return nil, gerror.New(consts.GetAuthErrorMessage(consts.ErrUserDisabled))
        default:
            return nil, gerror.New("用户状态异常，无法登录")
        }
    }

    // 获取用户角色信息
    userRole, err := s.getUserPrimaryRole(ctx, user.Id)
    if err != nil {
        return nil, gerror.Newf("获取用户角色失败: %v", err)
    }

    // 生成JWT Token
    payload := &simple.JWTPayload{
        UserId:   user.Id,
        Username: user.Username,
        RoleId:   userRole.RoleId,
        RoleKey:  userRole.RoleCode,
        DeptId:   user.DeptId,
        App:      consts.AppApi,
    }

    secretKey := simple.GetJWTSecretKey(ctx)
    accessToken, err := simple.GenerateJWTToken(payload, secretKey)
    if err != nil {
        return nil, gerror.Newf("生成访问令牌失败: %v", err)
    }

    // 生成刷新令牌
    refreshToken, err := s.generateRefreshToken(ctx, user.Id)
    if err != nil {
        return nil, gerror.Newf("生成刷新令牌失败: %v", err)
    }

    // 更新用户登录信息
    if err = s.updateLoginInfo(ctx, user.Id); err != nil {
        g.Log().Warningf(ctx, "更新用户登录信息失败: %v", err)
    }

    // 获取用户权限
    permissions, menuIds, err := s.getUserPermissions(ctx, user.Id)
    if err != nil {
        g.Log().Warningf(ctx, "获取用户权限失败: %v", err)
        permissions = []string{}
        menuIds = []int64{}
    }

    // 构建响应
    res = &sysout.LoginTokenModel{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        TokenType:    "Bearer",
        ExpiresIn:    24 * 3600, // 24小时
        UserInfo:     user,
        Permissions:  permissions,
        MenuIds:      menuIds,
    }

    return res, nil
}
```

### 模型层 (Model)

定义数据结构，包括请求参数、响应结果、数据库实体等。项目中实际包含以下子目块：

#### 实体模型 (Entity)
```go
// 用户实体（实际定义，337行）
type User struct {
    Id                   int64       `json:"id"                    description:"主键ID"`
    Username             string      `json:"username"              description:"用户名"`
    Email                string      `json:"email"                 description:"邮箱地址"`
    Phone                string      `json:"phone"                 description:"手机号码"`
    Password             string      `json:"-"                     description:"密码（不返回给前端）"`
    Salt                 string      `json:"-"                     description:"密码盐值（不返回给前端）"`
    RealName             string      `json:"realName"              description:"真实姓名"`
    Nickname             string      `json:"nickname"              description:"昵称"`
    Avatar               string      `json:"avatar"                description:"头像URL"`
    Gender               int         `json:"gender"                description:"性别"`
    Birthday             *gtime.Time `json:"birthday"              description:"生日"`
    DeptId               int64       `json:"deptId"                description:"部门ID"`
    Position             string      `json:"position"              description:"职位"`
    Status               int         `json:"status"                description:"状态"`
    LoginIp              string      `json:"loginIp"               description:"最后登录IP"`
    LoginAt              *gtime.Time `json:"loginAt"               description:"最后登录时间"`
    LoginCount           int         `json:"loginCount"            description:"登录次数"`
    PasswordResetToken   string      `json:"-"                     description:"密码重置令牌"`
    PasswordResetExpires *gtime.Time `json:"-"                     description:"密码重置过期时间"`
    EmailVerifiedAt      *gtime.Time `json:"emailVerifiedAt"       description:"邮箱验证时间"`
    PhoneVerifiedAt      *gtime.Time `json:"phoneVerifiedAt"       description:"手机验证时间"`
    TwoFactorEnabled     int         `json:"twoFactorEnabled"      description:"是否启用双因子认证"`
    TwoFactorSecret      string      `json:"-"                     description:"双因子认证密钥"`
    Remark               string      `json:"remark"                description:"备注说明"`
    CreatedBy            int64       `json:"createdBy"             description:"创建人ID"`
    UpdatedBy            int64       `json:"updatedBy"             description:"修改人ID"`
    CreatedAt            *gtime.Time `json:"createdAt"             description:"创建时间"`
    UpdatedAt            *gtime.Time `json:"updatedAt"             description:"更新时间"`
    DeletedAt            *gtime.Time `json:"deletedAt,omitempty"   description:"删除时间"`
}

// 角色实体（实际定义，180行）
type Role struct {
    Id          int64       `json:"id"          description:"主键ID"`
    Name        string      `json:"name"        description:"角色名称"`
    Code        string      `json:"code"        description:"角色编码"`
    Description string      `json:"description" description:"角色描述"`
    Status      int         `json:"status"      description:"状态：1=启用 0=禁用"`
    Sort        int         `json:"sort"        description:"排序号，数字越小越靠前"`
    DataScope   int         `json:"dataScope"   description:"数据权限范围"`
    Remark      string      `json:"remark"      description:"备注说明"`
    CreatedBy   int64       `json:"createdBy"   description:"创建人ID"`
    UpdatedBy   int64       `json:"updatedBy"   description:"修改人ID"`
    CreatedAt   *gtime.Time `json:"createdAt"   description:"创建时间"`
    UpdatedAt   *gtime.Time `json:"updatedAt"   description:"更新时间"`
}
```

#### 输入输出模型
- **输入模型 (input/sysin)**: 定义API请求参数结构
- **输出模型 (output/sysout)**: 定义API响应结果结构

### 数据层 (ORM)

通过GoFrame的ORM进行数据库操作，支持多种数据库：

```go
// 数据库操作示例
func (s *sUser) ValidateUser(ctx context.Context, username, password string) (user *sysout.UserModel, err error) {
    var userEntity *entity.User
    err = g.DB().Model("users").
        Where("username = ? AND deleted_at IS NULL", username).
        Scan(&userEntity)
    
    if err != nil {
        return nil, gerror.Newf("查询用户失败: %v", err)
    }

    // 验证密码逻辑...
    
    return sysout.ConvertToUserModel(userEntity), nil
}
```

## 分层设计原则

1. **单向依赖原则**: 上层可以调用下层，下层不能调用上层
2. **接口隔离原则**: Controller通过Service接口调用Logic，降低耦合
3. **关注点分离**: 每一层只关注自己的职责
4. **依赖注入**: 通过服务注册机制实现依赖注入

## 服务注册机制

项目采用服务注册机制来管理依赖关系：

```go
// 在逻辑层注册服务实现
func init() {
    service.RegisterUser(NewUser())
    service.RegisterRole(NewRole())
}

// 控制器通过服务接口调用
func (c *cUser) Login(ctx context.Context, req *user.LoginReq) (res *user.LoginRes, err error) {
    // 通过服务接口调用逻辑层实现
    out, err := service.User().Login(ctx, &req.UserLoginInp)
    if err != nil {
        return nil, err
    }
    return &user.LoginRes{LoginTokenModel: out}, nil
}
```

## 数据流转过程

一个完整的请求处理流程如下：

1. **HTTP 请求进入系统**
2. **中间件链处理**：认证、跨域、日志等通用逻辑
3. **路由层映射**：将请求映射到控制器方法
4. **控制器层处理**：
   - 接收请求参数并进行参数验证
   - 调用服务层接口
   - 处理返回结果并封装响应
5. **服务层接口调用**：通过接口定义调用逻辑层
6. **逻辑层业务处理**：
   - 实现具体的业务逻辑
   - 进行数据库操作
   - 处理业务规则和流程
7. **数据库操作**：通过ORM进行数据持久化
8. **响应返回**：逐层返回处理结果到客户端

## 用户登录模块调用示例

以用户登录为例，展示完整的调用流程：

```
HTTP POST /api/user/login
    ↓
Router → Controller.Login()
    ↓
service.User().Login() [接口调用]
    ↓
logic.sUser.Login() [实现调用]
    ↓
1. 验证用户密码 (数据库查询)
2. 检查用户状态
3. 获取用户角色信息
4. 生成JWT Token
5. 生成刷新令牌
6. 更新登录信息
7. 获取用户权限
    ↓
Return sysout.LoginTokenModel
    ↓
Controller 封装 user.LoginRes
    ↓
HTTP Response JSON
```

## 权限验证中间件示例

```go
// 权限验证中间件流程
func (s *sMiddleware) ApiAuth(r *ghttp.Request) {
    // 1. 获取Token
    token := s.getTokenFromRequest(r)
    
    // 2. 验证Token
    payload, err := simple.ParseJWTToken(token, secretKey)
    if err != nil {
        response.JsonExit(r, consts.ErrorTokenInvalid)
    }
    
    // 3. 验证用户状态
    user := s.getCurrentUser(r.Context(), payload.UserId)
    if !user.CanLogin() {
        response.JsonExit(r, consts.ErrorUserDisabled)
    }
    
    // 4. 验证权限
    if !s.hasPermission(r.Context(), payload.UserId, r.URL.Path) {
        response.JsonExit(r, consts.ErrorPermissionDenied)
    }
    
    // 5. 设置用户上下文
    r.SetCtxVar(consts.CtxUserId, payload.UserId)
    r.Middleware.Next()
}
```

## 总结

本项目的分层设计通过接口抽象实现了高内聚低耦合：

- **Router层**：专注于HTTP路由映射
- **Controller层**：专注于HTTP请求处理和响应封装
- **Service层**：提供业务接口抽象，实现依赖倒置
- **Logic层**：包含具体的业务逻辑实现，如用户认证（441行）、角色管理（1086行）
- **Model层**：定义数据结构和参数模型，包含完整的实体定义
- **ORM层**：提供数据库操作抽象

这种设计使得系统具有良好的扩展性、可测试性和可维护性，特别是通过服务注册机制实现了松耦合的依赖管理。
