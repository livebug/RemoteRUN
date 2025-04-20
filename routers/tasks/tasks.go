package tasks

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Router(e *gin.Engine) {
	// 任务路由
	task := e.Group("/tasks")
	{
		// 任务相关的路由
		task.GET("/getlogs", taskGetLogs)
		task.POST("/run", taskRun)
		// task.POST("/run_async", taskUpdate)
		// task.POST("/delete", taskDelete)
	}
}

// taskGetLogs handles the GET request for fetching task logs.
func taskGetLogs(c *gin.Context) {
	taskID := c.Query("taskid")
	if taskID == "" {
		c.JSON(400, gin.H{"error": "taskid is required"})
		return
	}

	logFilePath := "./tmp/tasks/logs/" + taskID + ".log"
	logFile, err := os.Open(logFilePath)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to open log file"})
		return
	}
	defer logFile.Close()

	// 读取日志文件内容
	content, err := io.ReadAll(logFile)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to read log file"})
		return
	}

	// 返回日志内容
	c.JSON(200, gin.H{
		"task_id": taskID,
		"logs":    string(content),
	})
	// 这里可以根据需要返回日志内容或其他信息
}

type Task struct {
	Command string   `json:"command"` // 命令
	Args    []string `json:"args"`    // 参数
	Order   int      `json:"order"`   // 执行顺序
}

func executeTask(tasks []Task, logFile *os.File) {

	fmt.Println("Executing tasks...")

	for _, task := range tasks {

		// 执行日志返回前端
		// 这里可以使用 exec.Command 来执行命令，并获取输出和错误信息
		cmd := exec.Command(task.Command, task.Args...)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			fmt.Println("Error creating stdout pipe:", err)
			return
		}

		stderr, err := cmd.StderrPipe()
		if err != nil {
			fmt.Println("Error creating stderr pipe:", err)
			return
		}

		// 启动命令
		if err := cmd.Start(); err != nil {
			fmt.Println("Error starting command:", err)
			return
		}

		go io.Copy(logFile, stdout) // 将命令的输出写入响应
		go io.Copy(logFile, stderr) // 将命令的错误输出写入响应

		if err := cmd.Wait(); err != nil {
			fmt.Println("Error waiting for command:", err)
			logFile.WriteString(fmt.Sprintf("Error executing command %s: %v\n", task.Command, err))
			return
		}
	}
	fmt.Println("All tasks executed successfully")
	logFile.Close()
}

// taskRun handles the POST request for running a task.
func taskRun(c *gin.Context) {
	// 读取任务组，一个json数组，里面多个task，task包含，命令command、参数args、执行顺序order
	tasks := []Task{}
	if err := c.ShouldBindJSON(&tasks); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	// 检查 order 是否重复
	// 创建一个 map 用于记录任务的执行顺序是否重复
	orderMap := make(map[int]bool)
	// 遍历任务列表中的每个任务
	for _, task := range tasks {
		// 检查当前任务的执行顺序是否已经存在于 map 中
		if _, exists := orderMap[task.Order]; exists {
			// 如果发现重复的执行顺序，返回 400 错误响应
			c.JSON(400, gin.H{"error": "Duplicate order found"})
			return
		}
		// 将当前任务的执行顺序标记为已存在
		orderMap[task.Order] = true
	}
	// 检查命令是否为空
	for _, task := range tasks {
		if task.Command == "" {
			c.JSON(400, gin.H{"error": "Command cannot be empty"})
			return
		}
	}

	// tasks 按顺序排序，当传入的任务组中有多个任务时，按顺序执行
	tasks = sortTasks(tasks)

	// 执行任务
	taskID := uuid.New().String()
	logFilePath := "./tmp/tasks/logs/" + taskID + ".log"

	// 创建日志文件
	if err := os.MkdirAll("./tmp/tasks/logs/", os.ModePerm); err != nil {
		fmt.Println("Failed to create log directory:", err)
		c.JSON(500, gin.H{"error": "Failed to create log directory"})
		return
	}
	// 创建日志文件
	logFile, err := os.Create(logFilePath)
	if err != nil {
		fmt.Println("Failed to create log file:", err)
		c.JSON(500, gin.H{"error": "Failed to create log file"})
		return
	}

	go executeTask(tasks, logFile)

	// 返回任务 ID 和日志文件路径
	c.JSON(200, gin.H{
		"task_id":       taskID,
		"log_file_path": logFilePath,
	})
}

func sortTasks(tasks []Task) []Task {
	// 按照 Order 字段排序任务
	// 这里可以使用 sort.Slice 来排序
	sort.Slice(tasks, func(i, j int) bool { return tasks[i].Order < tasks[j].Order })

	return tasks
}
