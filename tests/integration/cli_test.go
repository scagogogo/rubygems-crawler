package integration

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// CLI 集成测试需要构建二进制文件，所以这里采用调用外部命令的方式
func TestCLI(t *testing.T) {
	// 跳过长时间运行的测试
	if testing.Short() {
		t.Skip("在短模式下跳过CLI测试")
	}

	// 尝试获取编译后的二进制文件路径
	cmd := exec.Command("go", "build", "-o", "rubygems-cli", "../../cmd/rubygems/main.go")
	err := cmd.Run()
	if err != nil {
		t.Fatalf("编译CLI失败: %v", err)
	}

	// 测试帮助信息
	t.Run("显示帮助信息", func(t *testing.T) {
		output, err := exec.Command("./rubygems-cli", "-help").CombinedOutput()
		assert.NoError(t, err, "执行帮助命令失败")
		assert.Contains(t, string(output), "获取包信息", "帮助输出应包含功能描述")
		assert.Contains(t, string(output), "搜索包", "帮助输出应包含功能描述")
	})

	// 测试获取包信息
	t.Run("获取包信息", func(t *testing.T) {
		output, err := exec.Command("./rubygems-cli", "-get", "-gem", "rails").CombinedOutput()
		assert.NoError(t, err, "获取包信息失败")
		assert.Contains(t, string(output), "rails", "输出应包含包名")
	})

	// 测试搜索功能
	t.Run("搜索功能", func(t *testing.T) {
		output, err := exec.Command("./rubygems-cli", "-search", "-query", "rails", "-limit", "5").CombinedOutput()
		assert.NoError(t, err, "搜索包失败")
		assert.Contains(t, string(output), "rails", "搜索结果应包含rails")
	})

	// 测试获取版本信息
	t.Run("获取版本信息", func(t *testing.T) {
		output, err := exec.Command("./rubygems-cli", "-versions", "-gem", "rails", "-limit", "5").CombinedOutput()
		assert.NoError(t, err, "获取版本信息失败")
		assert.Contains(t, string(output), "rails", "版本信息应包含包名")
	})

	// 测试获取依赖信息
	t.Run("获取依赖信息", func(t *testing.T) {
		output, err := exec.Command("./rubygems-cli", "-deps", "-gem", "rails").CombinedOutput()
		assert.NoError(t, err, "获取依赖信息失败")
		assert.NotEmpty(t, output, "依赖信息不应为空")
	})

	// 测试获取反向依赖信息
	t.Run("获取反向依赖信息", func(t *testing.T) {
		output, err := exec.Command("./rubygems-cli", "-rdeps", "-gem", "rack", "-limit", "5").CombinedOutput()
		assert.NoError(t, err, "获取反向依赖信息失败")
		assert.NotEmpty(t, output, "反向依赖信息不应为空")
	})

	// 测试JSON输出
	t.Run("JSON输出", func(t *testing.T) {
		output, err := exec.Command("./rubygems-cli", "-get", "-gem", "rails", "-json").CombinedOutput()
		assert.NoError(t, err, "获取JSON格式的包信息失败")

		// 尝试解析JSON
		var result map[string]interface{}
		err = json.Unmarshal(output, &result)
		assert.NoError(t, err, "解析JSON输出失败")
		assert.Equal(t, "rails", result["name"], "JSON应包含正确的包名")
	})

	// 测试使用缓存
	t.Run("使用缓存", func(t *testing.T) {
		// 首次获取
		start := time.Now()
		_, err := exec.Command("./rubygems-cli", "-get", "-gem", "rails").CombinedOutput()
		assert.NoError(t, err, "首次获取包信息失败")
		firstDuration := time.Since(start)

		// 使用缓存再次获取
		start = time.Now()
		_, err = exec.Command("./rubygems-cli", "-get", "-gem", "rails", "-cache").CombinedOutput()
		assert.NoError(t, err, "使用缓存获取包信息失败")
		secondDuration := time.Since(start)

		// 缓存应该更快
		t.Logf("无缓存耗时: %v, 使用缓存耗时: %v", firstDuration, secondDuration)
	})

	// 测试镜像选择
	t.Run("镜像选择", func(t *testing.T) {
		mirrors := []string{"default", "ruby-china", "tsinghua", "aliyun"}

		for _, mirror := range mirrors {
			t.Run(mirror, func(t *testing.T) {
				output, err := exec.Command("./rubygems-cli", "-get", "-gem", "rake", "-mirror", mirror).CombinedOutput()
				assert.NoError(t, err, "使用镜像 %s 获取包信息失败", mirror)
				assert.Contains(t, string(output), "rake", "使用镜像 %s 的输出应包含包名", mirror)
			})
		}
	})

	// 测试无效的命令
	t.Run("无效的命令", func(t *testing.T) {
		cmd := exec.Command("./rubygems-cli", "-invalid", "-gem", "rails")
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		err := cmd.Run()

		// 命令应该退出且有错误信息
		assert.Error(t, err, "无效命令应返回错误")
	})

	// 测试清理
	defer func() {
		err := exec.Command("rm", "-f", "rubygems-cli").Run()
		if err != nil {
			t.Logf("清理文件失败: %v", err)
		}
	}()
}
