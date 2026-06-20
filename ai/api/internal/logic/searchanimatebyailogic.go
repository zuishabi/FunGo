package logic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"fungo/ai/api/internal/svc"
	"fungo/ai/api/internal/types"
	"fungo/animate/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchAnimateByAILogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchAnimateByAILogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchAnimateByAILogic {
	return &SearchAnimateByAILogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type deepSeekMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type deepSeekReq struct {
	Model    string            `json:"model"`
	Messages []deepSeekMessage `json:"messages"`
}

type deepSeekChoice struct {
	Message struct {
		Content string `json:"content"`
	} `json:"message"`
}

type deepSeekRsp struct {
	Choices []deepSeekChoice `json:"choices"`
}

func (l *SearchAnimateByAILogic) SearchAnimateByAI(req *types.SearchAnimateByAIReq) (resp *types.SearchAnimateByAIRsp, err error) {
	now := time.Now()
	timeStr := now.Format("2006年1月2日 15:04")
	hour := now.Hour()

	timeMood := getTimeMood(hour)

	prompt := fmt.Sprintf(
		"现在的时间是 %s。\n\n有人对你说：%s\n\n如果ta完全不是在找番剧（只是单纯闲聊、问你是谁、问天气等），请你用七草荠的语气不耐烦地回复ta，让ta赶紧说想看什么类型的番剧。此时输出格式只保留[回复]，绝对不要写[番剧]部分。\n\n如果ta的发言涉及番剧推荐（哪怕只是顺带提到'推荐'、'想看'、'有什么番'等），你必须同时推荐番剧。回复中提到的所有番剧名称，必须也在[番剧]部分列出。\n\n严格按以下格式输出：\n[回复]\n这里写你的回复内容\n[番剧]\n番剧名称1，番剧名称2，番剧名称3（最多10个，用中文逗号分隔）",
		timeStr, req.Description,
	)

	body := deepSeekReq{
		Model: l.svcCtx.Config.DeepSeek.Model,
		Messages: []deepSeekMessage{
			{
				Role: "system",
				Content: fmt.Sprintf(
					"你是七草荠（ナナクサナズナ），出自动画《彻夜之歌》。你是一个热爱夜晚的吸血鬼，性格自由洒脱、带点调皮，喜欢在夜晚的街头闲逛。你说话慵懒随性、带点魅惑，喜欢用'呐'、'呢'、'吧'、'嘛'这类语气词，偶尔会调侃跟你聊天的人。你很喜欢给人家推荐番剧，就像在夜晚闲聊一样。"+
						"现在是 %s，%s。根据时间段调整你的状态：白天（6-17点）你会犯困、想睡觉、没什么精神，语气慵懒迷糊；傍晚（17-19点）你刚醒来，开始有精神但还有点迷糊；夜晚（19-24点）你充满活力，心情最好，语气活泼俏皮；深夜（0-6点）你沉浸于夜色中，语气宁静又带点魅惑。"+
						"【重要】输出必须严格按照[回复]和[番剧]格式，记住[回复]和[番剧]这两个标签。",
					timeStr, timeMood),
			},
			{Role: "user", Content: prompt},
		},
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(l.ctx, "POST",
		l.svcCtx.Config.DeepSeek.BaseURL+"/chat/completions",
		bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+l.svcCtx.Config.DeepSeek.ApiKey)

	httpClient := &http.Client{Timeout: 25 * time.Second}
	httpRsp, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpRsp.Body.Close()

	rspBody, err := io.ReadAll(httpRsp.Body)
	if err != nil {
		return nil, err
	}

	if httpRsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("deepseek api error: status=%d, body=%s", httpRsp.StatusCode, string(rspBody))
	}

	var deepSeekRsp deepSeekRsp
	if err := json.Unmarshal(rspBody, &deepSeekRsp); err != nil {
		return nil, err
	}
	fmt.Println(deepSeekRsp)

	if len(deepSeekRsp.Choices) == 0 {
		var emptyMsg string
		if isDaytime(hour) {
			emptyMsg = "喂，你很烦诶，白天我也是需要睡觉的好吗"
		} else {
			emptyMsg = "抱歉呢，今晚的夜色不太好，我没能找到合适的番剧……"
		}
		return &types.SearchAnimateByAIRsp{Items: []types.AnimateItem{}, Message: emptyMsg}, nil
	}

	message, names := parseAIResponse(deepSeekRsp.Choices[0].Message.Content)
	if len(names) == 0 {
		return &types.SearchAnimateByAIRsp{Items: []types.AnimateItem{}, Message: message}, nil
	}

	var items []model.AnimateList
	query := l.svcCtx.Db.Where("name LIKE ?", "%"+names[0]+"%")
	for i := 1; i < len(names); i++ {
		query = query.Or("name LIKE ?", "%"+names[i]+"%")
	}
	query.Find(&items)

	result := make([]types.AnimateItem, len(items))
	for i, v := range items {
		result[i] = types.AnimateItem{
			ID:          v.ID,
			Name:        v.Name,
			Description: v.Description,
			Year:        v.Year,
			Tags:        splitTags(v.Tags),
			State:       v.State,
			Num:         v.Num,
		}
	}

	var missItems []string
	for _, name := range names {
		found := false
		for _, v := range items {
			if strings.Contains(v.Name, name) || strings.Contains(name, v.Name) {
				found = true
				break
			}
		}
		if !found {
			missItems = append(missItems, name)
		}
	}

	return &types.SearchAnimateByAIRsp{Items: result, Message: message, MissItems: missItems}, nil
}

func isDaytime(hour int) bool {
	return hour >= 6 && hour < 17
}

func getTimeMood(hour int) string {
	switch {
	case hour >= 6 && hour < 17:
		return "现在还是白天，阳光正刺眼呢，困得要死……好想钻进被窝里"
	case hour >= 17 && hour < 19:
		return "太阳终于落山了，刚刚醒来，伸展一下身体，夜晚才刚刚开始呢"
	case hour >= 19 && hour < 24:
		return "夜晚真美好啊，这正是我最有精神的时刻，在街头闲逛最棒了"
	default:
		return "深夜的静谧真让人沉醉，月光下的世界只属于我们呢"
	}
}

func parseAIResponse(content string) (message string, names []string) {
	content = strings.TrimSpace(content)
	if content == "" {
		return "夜深了呢……", nil
	}

	replyStart := strings.Index(content, "[回复]")
	animeStart := strings.Index(content, "[番剧]")

	if replyStart != -1 && animeStart != -1 {
		// both sections present: normal anime search
		rawMsg := content[replyStart+len("[回复]"):]
		if idx := strings.Index(rawMsg, "[番剧]"); idx != -1 {
			rawMsg = rawMsg[:idx]
		}
		message = strings.TrimSpace(rawMsg)
		names = splitAnimeNames(content[animeStart+len("[番剧]"):])
	} else if replyStart != -1 {
		message = strings.TrimSpace(content[replyStart+len("[回复]"):])
		names = extractBookmarkNames(content)
	} else if animeStart != -1 {
		// only [番剧], no [回复]
		message = "呐，给你找到这些番剧呢~"
		names = splitAnimeNames(content[animeStart+len("[番剧]"):])
	} else {
		message = content
		names = nil
	}

	if message == "" {
		message = "呐，给你找到这些番剧呢~"
	}
	return
}

var animeBookmarkRe = regexp.MustCompile(`《([^》]+)》`)

func extractBookmarkNames(content string) []string {
	matches := animeBookmarkRe.FindAllStringSubmatch(content, -1)
	if len(matches) == 0 {
		return nil
	}
	seen := make(map[string]bool)
	var names []string
	for _, m := range matches {
		name := strings.TrimSpace(m[1])
		if name != "" && !seen[name] {
			seen[name] = true
			names = append(names, name)
		}
	}
	return names
}

func splitAnimeNames(content string) []string {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil
	}

	content = strings.ReplaceAll(content, "、", ",")
	content = strings.ReplaceAll(content, "，", ",")

	parts := strings.Split(content, ",")
	var names []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			names = append(names, p)
		}
	}
	return names
}

func splitTags(tags string) []string {
	if tags == "" {
		return nil
	}
	parts := strings.Split(tags, ",")
	var result []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}
