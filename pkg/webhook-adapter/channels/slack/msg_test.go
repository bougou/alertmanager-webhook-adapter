package slack

import (
	"fmt"
	"os"
	"testing"

	"github.com/slack-go/slack"
)

func Test_SlackSender2(t *testing.T) {
	token := os.Getenv("SLACK_APP_TOKEN")
	channel := "#jenkins-ci"
	fmt.Println("slack token:", token)
	sender := NewSender(token, channel, MsgTypeMarkdown)
	blocks := createSlackBlocks()
	msg := Msg(blocks)
	if err := sender.SendMsg(msg); err != nil {
		t.Error(err)
	}
}

// slack.Block(s) can be used saas slack.MsgOption(s)
func createSlackBlocks() []slack.Block {
	// Header Section
	headerText := slack.NewTextBlockObject("mrkdwn", "*Where should we order lunch from?* Poll by <fakeLink.toUser.com|Mark>", false, false)
	headerSection := slack.NewSectionBlock(headerText, nil, nil)

	// Header2
	headerText2 := slack.NewHeaderBlock( // Header Block must be plain-text
		slack.NewTextBlockObject(
			slack.PlainTextType,
			"Hello world:fire:",
			true,  // only plain-text can set to true
			false, // must set to false for plain-text
		),
	)

	// Divider
	divSection := slack.NewDividerBlock()

	// SectionBlock
	textSection := slack.NewSectionBlock(
		slack.NewTextBlockObject(
			slack.PlainTextType,
			"### Hello\nabc",
			true,
			false,
		),
		[]*slack.TextBlockObject{
			slack.NewTextBlockObject(
				slack.MarkdownType,
				"test\n*HELO*\n> test",
				false, // must be set to false for markdown type
				true,
			),
		},
		nil,
	)

	// Shared Assets for example
	voteBtnText := slack.NewTextBlockObject("plain_text", "Vote", true, false)
	voteBtnEle := slack.NewButtonBlockElement("", "click_me_123", voteBtnText)
	profileOne := slack.NewImageBlockElement("https://api.slack.com/img/blocks/bkb_template_images/profile_1.png", "Michael Scott")
	profileTwo := slack.NewImageBlockElement("https://api.slack.com/img/blocks/bkb_template_images/profile_2.png", "Dwight Schrute")
	profileThree := slack.NewImageBlockElement("https://api.slack.com/img/blocks/bkb_template_images/profile_3.png", "Pam Beasely")
	profileFour := slack.NewImageBlockElement("https://api.slack.com/img/blocks/bkb_template_images/profile_4.png", "Angela")

	// Option One Info
	optOneText := slack.NewTextBlockObject("mrkdwn", ":sushi: *Ace Wasabi Rock-n-Roll Sushi Bar*\nThe best landlocked sushi restaurant.", false, false)
	optOneSection := slack.NewSectionBlock(optOneText, nil, slack.NewAccessory(voteBtnEle))

	// Option One Votes
	optOneVoteText := slack.NewTextBlockObject("plain_text", "3 votes", true, false)
	optOneContext := slack.NewContextBlock("", []slack.MixedElement{profileOne, profileTwo, profileThree, optOneVoteText}...)

	// Option Two Info
	optTwoText := slack.NewTextBlockObject("mrkdwn", ":hamburger: *Super Hungryman Hamburgers*\nOnly for the hungriest of the hungry.", false, false)
	optTwoSection := slack.NewSectionBlock(optTwoText, nil, slack.NewAccessory(voteBtnEle))

	// Option Two Votes
	optTwoVoteText := slack.NewTextBlockObject("plain_text", "2 votes", true, false)
	optTwoContext := slack.NewContextBlock("", []slack.MixedElement{profileFour, profileTwo, optTwoVoteText}...)

	// Option Three Info
	optThreeText := slack.NewTextBlockObject("mrkdwn", ":ramen: *Kagawa-Ya Udon Noodle Shop*\nDo you like to shop for noodles? We have noodles.", false, false)
	optThreeSection := slack.NewSectionBlock(optThreeText, nil, slack.NewAccessory(voteBtnEle))

	// Option Three Votes
	optThreeVoteText := slack.NewTextBlockObject("plain_text", "No votes", true, false)
	optThreeContext := slack.NewContextBlock("", []slack.MixedElement{optThreeVoteText}...)

	// Suggestions Action
	btnTxt := slack.NewTextBlockObject("plain_text", "Add a suggestion", false, false)
	nextBtn := slack.NewButtonBlockElement("", "click_me_123", btnTxt)
	actionBlock := slack.NewActionBlock("", nextBtn)

	blocks := []slack.Block{
		headerSection,
		divSection,
		headerText2,
		textSection,
		optOneSection,
		optOneContext,
		optTwoSection,
		optTwoContext,
		optThreeSection,
		optThreeContext,
		divSection,
		actionBlock,
	}

	return blocks
}
