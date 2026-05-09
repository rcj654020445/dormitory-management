# 测试标准

## 1 测试组织

```
backend/
├── internal/
│   ├── service/
│   │   └── student_svc_test.go     # Service 单元测试
│   └── repository/
│       └── student_repo_test.go    # Repository 测试（使用 mock DB）
└── cmd/
    └── server_test.go               # 集成测试

frontend/
├── src/
│   ├── stores/
│   │   └── student.test.ts         # Store 测试
│   └── views/
│       └── StudentList.test.vue    # 组件测试
└── e2e/
    └── student.spec.ts              # 端到端测试
```

## 2 运行测试

### Backend (Go)

```bash
# 所有测试
go test ./...

# 带覆盖率
go test -cover ./...

# 指定包
go test ./internal/service/...

# 详细输出
go test -v ./...
```

### Frontend (Vue3)

```bash
cd frontend
npm run test          # 运行 Vitest
npm run test:unit     # 仅运行单元测试
npm run test:e2e      # E2E 测试 (Cypress/Playwright)
```

## 3 测试模式

### 3.1 单元测试 (Go)

```go
func TestStudentService_CreateStudent(t *testing.T) {
    repo := &mockStudentRepository{}
    svc := service.NewStudentService(repo)

    req := &types.CreateStudentRequest{
        StudentID: "2024001",
        Name:      "张三",
        Gender:    "male",
        Phone:     "13800138000",
        Email:     "zhangsan@example.com",
        Major:     "计算机科学与技术",
        Grade:     1,
    }

    student, err := svc.CreateStudent(context.Background(), req)

    require.NoError(t, err)
    require.NotNil(t, student)
    assert.Equal(t, "2024001", student.StudentID)
    assert.Equal(t, "pending", student.Status)
}
```

### 3.2 表格驱动测试 (Go)

```go
func TestAllocation_Validate(t *testing.T) {
    tests := []struct {
        name    string
        student  *types.Student
        room    *types.Room
        wantErr bool
    }{
        {
            name:    "gender mismatch",
            student: &types.Student{Gender: "male"},
            room:    &types.Room{BuildingGender: "female"},
            wantErr: true,
        },
        {
            name:    "room full",
            student: &types.Student{Gender: "male"},
            room:    &types.Room{BedsUsed: 4, BedsTotal: 4},
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validateAllocation(tt.student, tt.room)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

## 4 覆盖率目标

| Component | Target | Current |
|-----------|--------|---------|
| Service layer | 80% | N/A |
| Repository layer | 70% | N/A |
| HTTP handlers | 60% | N/A |
| Core utilities | 90% | N/A |

## 5 Mocking

使用接口来管理依赖：

```go
type StudentRepository interface {
    Create(ctx context.Context, student *types.Student) error
    GetByID(ctx context.Context, id string) (*types.Student, error)
    // ...
}
```

用于测试的 Mock 实现：

```go
type mockStudentRepository struct {
    students map[string]*types.Student
}

func (m *mockStudentRepository) Create(ctx context.Context, student *types.Student) error {
    m.students[student.ID] = student
    return nil
}
```