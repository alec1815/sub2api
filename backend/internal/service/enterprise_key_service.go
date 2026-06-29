package service

import (
	"context"
	"fmt"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
)

// Enterprise key purpose and tool constants
const (
	ToolCursor     = "cursor"
	ToolTrae       = "trae"
	ToolClaudeCode = "claude_code"
	ToolCodex      = "codex"
	ToolOpenCode   = "opencode"
	ToolPixso      = "pixso"
	ToolOther      = "other"

	MaxUsagePurposeLen = 200
)

var (
	ErrKeyUsagePurposeTooLong = infraerrors.BadRequest("USAGE_PURPOSE_TOO_LONG", fmt.Sprintf("usage purpose must be %d characters or less", MaxUsagePurposeLen))
	ErrKeyBoundToolInvalid    = infraerrors.BadRequest("BOUND_TOOL_INVALID", "bound tool must be one of: cursor, trae, claude_code, codex, opencode, pixso, other")
)

// validBoundTools is the set of allowed bound_tool values.
var validBoundTools = map[string]bool{
	ToolCursor:     true,
	ToolTrae:       true,
	ToolClaudeCode: true,
	ToolCodex:      true,
	ToolOpenCode:   true,
	ToolPixso:      true,
	ToolOther:      true,
}

// CreateEnterpriseKeyRequest 创建企业密钥请求
type CreateEnterpriseKeyRequest struct {
	Name         string  `json:"name"`
	GroupIDs     []int64 `json:"group_ids"`
	AssignedTo   *int64  `json:"assigned_to"`   // 分配给成员
	UsagePurpose string  `json:"usage_purpose"` // 用途描述
	BoundTool    string  `json:"bound_tool"`    // 绑定的工具
}

// UpdateEnterpriseKeyRequest 更新企业密钥
type UpdateEnterpriseKeyRequest struct {
	Name         *string  `json:"name"`
	AssignedTo   *int64   `json:"assigned_to"`
	UsagePurpose *string  `json:"usage_purpose"`
	BoundTool    *string  `json:"bound_tool"`
	GroupIDs     []int64  `json:"group_ids"`
}

// EnterpriseKeyService 企业密钥管理
type EnterpriseKeyService struct {
	apiKeyRepo APIKeyRepository
	memberRepo EnterpriseMemberRepository
	entRepo    EnterpriseRepository
	groupRepo  APIKeyGroupRepository
}

// NewEnterpriseKeyService 构造函数
func NewEnterpriseKeyService(
	apiKeyRepo APIKeyRepository,
	memberRepo EnterpriseMemberRepository,
	entRepo EnterpriseRepository,
	groupRepo APIKeyGroupRepository,
) *EnterpriseKeyService {
	return &EnterpriseKeyService{
		apiKeyRepo: apiKeyRepo,
		memberRepo: memberRepo,
		entRepo:    entRepo,
		groupRepo:  groupRepo,
	}
}

// CreateEnterpriseKey 创建企业密钥
func (s *EnterpriseKeyService) CreateEnterpriseKey(ctx context.Context, adminUserID int64, req CreateEnterpriseKeyRequest) (*APIKey, error) {
	// 1. 校验用途字段
	if len(req.UsagePurpose) > MaxUsagePurposeLen {
		return nil, ErrKeyUsagePurposeTooLong
	}
	if req.BoundTool != "" && !validBoundTools[req.BoundTool] {
		return nil, ErrKeyBoundToolInvalid
	}

	// 2. 校验管理员有企业管理员身份
	member, err := s.memberRepo.GetByUserID(ctx, adminUserID)
	if err != nil {
		return nil, fmt.Errorf("get member: %w", err)
	}
	if member.Role != EnterpriseMemberRoleAdmin {
		return nil, infraerrors.Forbidden("ENTERPRISE_ADMIN_REQUIRED", "enterprise admin permission required")
	}
	if member.Status != StatusActive {
		return nil, ErrMemberEnterpriseNotActive
	}

	// 3. 校验企业 active
	ent, err := s.entRepo.GetByID(ctx, member.EnterpriseID)
	if err != nil {
		return nil, fmt.Errorf("get enterprise: %w", err)
	}
	if ent.Status != EnterpriseStatusActive {
		return nil, ErrMemberEnterpriseNotActive
	}

	// 4. 如果 assigned_to 不为空，校验目标成员属于同一企业
	if req.AssignedTo != nil && *req.AssignedTo > 0 {
		targetMember, memberErr := s.memberRepo.GetByID(ctx, *req.AssignedTo)
		if memberErr != nil {
			return nil, fmt.Errorf("get target member: %w", memberErr)
		}
		if targetMember.EnterpriseID != member.EnterpriseID {
			return nil, infraerrors.BadRequest("MEMBER_NOT_IN_ENTERPRISE", "target member does not belong to this enterprise")
		}
	}

	// 5. 创建 api_key（通过 APIKeyRepository.Create）
	// assignedTo 字段由 APIKey 的 AssignedTo 字段承载
	key := &APIKey{
		UserID:       adminUserID,
		Name:         req.Name,
		AssignedTo:   req.AssignedTo,
		UsagePurpose: req.UsagePurpose,
		BoundTool:    req.BoundTool,
		Status:       StatusActive,
	}
	if err := s.apiKeyRepo.Create(ctx, key); err != nil {
		return nil, fmt.Errorf("create api key: %w", err)
	}

	// 6. 设置分组关联
	if len(req.GroupIDs) > 0 {
		if err := s.groupRepo.SetGroups(ctx, key.ID, req.GroupIDs); err != nil {
			return nil, fmt.Errorf("set key groups: %w", err)
		}
	}

	return key, nil
}

// UpdateEnterpriseKey 更新企业密钥信息
func (s *EnterpriseKeyService) UpdateEnterpriseKey(ctx context.Context, id int64, req UpdateEnterpriseKeyRequest) (*APIKey, error) {
	key, err := s.apiKeyRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get key: %w", err)
	}

	if req.Name != nil {
		key.Name = *req.Name
	}
	if req.AssignedTo != nil {
		key.AssignedTo = req.AssignedTo
	}
	if req.UsagePurpose != nil {
		if len(*req.UsagePurpose) > MaxUsagePurposeLen {
			return nil, ErrKeyUsagePurposeTooLong
		}
		key.UsagePurpose = *req.UsagePurpose
	}
	if req.BoundTool != nil {
		if *req.BoundTool != "" && !validBoundTools[*req.BoundTool] {
			return nil, ErrKeyBoundToolInvalid
		}
		key.BoundTool = *req.BoundTool
	}

	if err := s.apiKeyRepo.Update(ctx, key); err != nil {
		return nil, fmt.Errorf("update key: %w", err)
	}

	if req.GroupIDs != nil {
		if err := s.groupRepo.SetGroups(ctx, id, req.GroupIDs); err != nil {
			return nil, fmt.Errorf("set key groups: %w", err)
		}
	}

	return key, nil
}

// ListEnterpriseKeys 企业密钥列表（跨 user_id + assigned_to 的 UNION 查询）
// NOTE: This requires the APIKeyRepository to support enterprise listing.
// Implementation detail will be resolved in the repository layer.
func (s *EnterpriseKeyService) ListEnterpriseKeys(
	ctx context.Context,
	enterpriseID int64,
	params pagination.PaginationParams,
) ([]APIKey, *pagination.PaginationResult, error) {
	// 获取企业所有 active 成员
	members, _, err := s.memberRepo.ListByEnterprise(ctx, enterpriseID,
		pagination.DefaultPagination(),
		EnterpriseMemberListFilters{Status: StatusActive},
	)
	if err != nil {
		return nil, nil, fmt.Errorf("list members: %w", err)
	}

	// 收集所有成员 user_id 和 member_id
	var userIDs []int64
	for _, m := range members {
		userIDs = append(userIDs, m.UserID)
	}

	// 通过现有 API 按 user 查 key（简化版，实际需 UNION 优化）
	// TODO: 实现 UNION 查询替代此逐用户查询
	var allKeys []APIKey
	for _, uid := range userIDs {
		keys, _, kerr := s.apiKeyRepo.ListByUserID(ctx, uid, pagination.DefaultPagination(), APIKeyListFilters{})
		if kerr != nil {
			continue
		}
		allKeys = append(allKeys, keys...)
	}

	// 简易分页
	total := int64(len(allKeys))
	offset := int64(params.Offset())
	limit := int64(params.Limit())
	if offset >= total {
		return nil, simplePaginationResult(total, params), nil
	}
	end := offset + limit
	if end > total {
		end = total
	}
	page := allKeys[offset:end]
	return page, simplePaginationResult(total, params), nil
}

// ToggleEnterpriseKey 启停企业密钥（active ↔ disabled）
func (s *EnterpriseKeyService) ToggleEnterpriseKey(ctx context.Context, adminUserID int64, keyID int64) (*APIKey, error) {
	key, err := s.apiKeyRepo.GetByID(ctx, keyID)
	if err != nil {
		return nil, fmt.Errorf("get key: %w", err)
	}

	// 校验操作者是该企业的管理员
	member, err := s.memberRepo.GetByUserID(ctx, adminUserID)
	if err != nil {
		return nil, fmt.Errorf("get member: %w", err)
	}
	if member.Role != EnterpriseMemberRoleAdmin {
		return nil, infraerrors.Forbidden("ENTERPRISE_ADMIN_REQUIRED", "enterprise admin permission required")
	}

	// 校验 key 属于该企业
	if key.AssignedTo == nil {
		if key.UserID != adminUserID {
			return nil, infraerrors.Forbidden("KEY_NOT_IN_ENTERPRISE", "key does not belong to your enterprise")
		}
	} else {
		targetMember, memberErr := s.memberRepo.GetByID(ctx, *key.AssignedTo)
		if memberErr != nil {
			return nil, fmt.Errorf("get target member: %w", memberErr)
		}
		if targetMember.EnterpriseID != member.EnterpriseID {
			return nil, infraerrors.Forbidden("KEY_NOT_IN_ENTERPRISE", "key does not belong to your enterprise")
		}
	}

	// 切换状态
	if key.Status == StatusActive {
		key.Status = StatusDisabled
	} else {
		key.Status = StatusActive
	}

	if err := s.apiKeyRepo.Update(ctx, key); err != nil {
		return nil, fmt.Errorf("update key status: %w", err)
	}

	return key, nil
}

// DeleteEnterpriseKeyWrapper 删除企业密钥（带权限校验）
func (s *EnterpriseKeyService) DeleteEnterpriseKeyWrapper(ctx context.Context, adminUserID int64, keyID int64) error {
	key, err := s.apiKeyRepo.GetByID(ctx, keyID)
	if err != nil {
		return fmt.Errorf("get key: %w", err)
	}

	// 校验操作者是该企业的管理员
	member, err := s.memberRepo.GetByUserID(ctx, adminUserID)
	if err != nil {
		return fmt.Errorf("get member: %w", err)
	}
	if member.Role != EnterpriseMemberRoleAdmin {
		return infraerrors.Forbidden("ENTERPRISE_ADMIN_REQUIRED", "enterprise admin permission required")
	}

	// 校验 key 属于该企业
	if key.AssignedTo == nil {
		// 管理员自用 key → user_id 就是管理员
		if key.UserID != adminUserID {
			return infraerrors.Forbidden("KEY_NOT_IN_ENTERPRISE", "key does not belong to your enterprise")
		}
	} else {
		targetMember, memberErr := s.memberRepo.GetByID(ctx, *key.AssignedTo)
		if memberErr != nil {
			return fmt.Errorf("get target member: %w", memberErr)
		}
		if targetMember.EnterpriseID != member.EnterpriseID {
			return infraerrors.Forbidden("KEY_NOT_IN_ENTERPRISE", "key does not belong to your enterprise")
		}
	}

		return s.apiKeyRepo.Delete(ctx, keyID)
}

// simplePaginationResult builds a PaginationResult from total and params.
func simplePaginationResult(total int64, params pagination.PaginationParams) *pagination.PaginationResult {
	pageSize := params.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	pages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		pages++
	}
	return &pagination.PaginationResult{
		Total:    total,
		Page:     params.Page,
		PageSize: pageSize,
		Pages:    pages,
	}
}
