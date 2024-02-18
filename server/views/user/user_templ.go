// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.543
package user

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import (
	"github.com/diegoandino/wonder-go/model"
	"html/template"
)

func Show(u model.UserPayload, friends []model.UserPayload) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<!doctype html><html class=\"scroll-smooth\"><head><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><script src=\"https://unpkg.com/htmx.org\"></script><script src=\"https://kit.fontawesome.com/166d97ea2f.js\" crossorigin=\"anonymous\"></script><script src=\"https://cdnjs.cloudflare.com/ajax/libs/flowbite/2.3.0/flowbite.min.js\"></script><link href=\"/static/styles.css\" rel=\"stylesheet\"><link href=\"https://cdnjs.cloudflare.com/ajax/libs/flowbite/2.3.0/flowbite.min.css\" rel=\"stylesheet\"></head>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = Navbar().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"pt-10\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = CurrentUser(u).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = Friends(friends).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}

func Navbar() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var2 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var2 == nil {
			templ_7745c5c3_Var2 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<nav class=\"border-gray-200 dark:bg-gray-900 z-50 bg-transparent backdrop-filter backdrop-blur-xl\"><div class=\"max-w-screen-xl flex flex-wrap items-center justify-between mx-auto p-4\"><a href=\"/home\" class=\"flex items-center space-x-3 rtl:space-x-reverse\"><img src=\"https://flowbite.com/docs/images/logo.svg\" class=\"h-8\" alt=\"Flowbite Logo\"> <span class=\"self-center text-2xl font-semibold whitespace-nowrap dark:text-white\">Wonder</span></a><div class=\"flex md:order-2\"><button type=\"button\" data-collapse-toggle=\"navbar-search\" aria-controls=\"navbar-search\" aria-expanded=\"false\" class=\"md:hidden text-gray-500 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700 focus:outline-none focus:ring-4 focus:ring-gray-200 dark:focus:ring-gray-700 rounded-lg text-sm p-2.5 me-1\"><svg class=\"w-5 h-5\" aria-hidden=\"true\" xmlns=\"http://www.w3.org/2000/svg\" fill=\"none\" viewBox=\"0 0 20 20\"><path stroke=\"currentColor\" stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"m19 19-4-4m0-7A7 7 0 1 1 1 8a7 7 0 0 1 14 0Z\"></path></svg> <span class=\"sr-only\">Search</span></button><div class=\"relative hidden md:block\"><div class=\"absolute inset-y-0 start-0 flex items-center ps-3 pointer-events-none\"><svg class=\"w-4 h-4 text-gray-500 dark:text-gray-400\" aria-hidden=\"true\" xmlns=\"http://www.w3.org/2000/svg\" fill=\"none\" viewBox=\"0 0 20 20\"><path stroke=\"currentColor\" stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"m19 19-4-4m0-7A7 7 0 1 1 1 8a7 7 0 0 1 14 0Z\"></path></svg> <span class=\"sr-only\">Search icon</span></div><input type=\"text\" id=\"search-navbar\" class=\"block w-full p-2 ps-10 text-sm text-gray-900 border border-gray-300 rounded-lg bg-gray-50 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500\" name=\"query\" hx-get=\"/search-friends\" hx-target=\"#search-results-dropdown\" hx-trigger=\"keyup\" hx-indicator=\"#loading-indicator\" hx-empty=\"document.getElementById(&#39;search-results-dropdown&#39;).style.display=&#39;none&#39;;\" placeholder=\"Search...\" oninput=\"hideIfEmpty()\"></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = NotificationBell().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div><div class=\"items-center justify-between hidden w-full md:flex md:w-auto md:order-1\" id=\"navbar-search\"><div class=\"relative mt-3 md:hidden\"><div class=\"absolute inset-y-0 start-0 flex items-center ps-3 pointer-events-none\"><svg class=\"w-4 h-4 text-gray-500 dark:text-gray-400\" aria-hidden=\"true\" xmlns=\"http://www.w3.org/2000/svg\" fill=\"none\" viewBox=\"0 0 20 20\"><path stroke=\"currentColor\" stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"m19 19-4-4m0-7A7 7 0 1 1 1 8a7 7 0 0 1 14 0Z\"></path></svg></div><input type=\"text\" id=\"search-navbar\" class=\"block w-full p-2 ps-10 text-sm text-gray-900 border border-gray-300 rounded-lg bg-gray-50 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500\" name=\"query\" hx-get=\"/search-friends\" hx-target=\"#search-results-dropdown\" hx-trigger=\"keyup\" placeholder=\"Search users...\" oninput=\"hideIfEmpty()\"></div><div id=\"search-results-dropdown\" class=\"hidden\" hx-swap-oob=\"delete\"></div></div></div></nav><script>\n\t\tfunction hideIfEmpty() {\n\t\t\tvar input = document.getElementById('search-navbar').value;\n\t\t\tvar inputLength = input.length;\n\t\t\tconsole.log(input);\n\t\t\tif (inputLength === 0) {\n\t\t\t\tdocument.getElementById('search-results-dropdown').style.display = 'none';\n\t\t\t}\n\t\t}\n\t</script>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}

const currentUserTemplate = `
    <div hx-get="/get-user-payload" hx-trigger="every 5s" hx-swap="outerHTML" class="my-4">
        <div class="currently-playing mt-4 mr-4 ml-4 p-4 rounded-lg shadow-md relative overflow-hidden">
            <div class="absolute inset-0 -z-10 bg-cover bg-center blur-xl" style="background-image: url('{{.CurrentAlbumArt}}');"></div>
            <div class="user-profile flex mb-2 z-10 relative">
                <img src="{{.ProfilePicture}}" alt="Profile Picture" class="w-16 h-16 rounded-full mr-4"/>
                <div>
                    <h3 class="nunito-bold text-xl text-white">{{.Username}}</h3>
                    <p class="text-sm text-white nunito-medium-italic">Currently Playing</p>
                </div>
            </div>
            <div class="ml-20 z-10 relative">
                <img src="{{.CurrentAlbumArt}}" alt="Album Art" class="w-36 h-36 mb-2"/>
                <div class="bg-black bg-opacity-20 backdrop-blur-lg rounded w-36">
                    <p class="text-md text-center text-white nunito-bold-italic">{{.CurrentSongName}}</p>
                    <p class="text-sm text-center text-white nunito-medium">by {{.CurrentArtistName}}</p>
                    <p class="text-sm text-center text-white nunito-semibold">{{.CurrentAlbumName}}</p>
                </div>
            </div>
        </div>
    </div>
`

var parsedCurrentUserTemplate = template.Must(template.New("current_user_template").Parse(currentUserTemplate))

func CurrentUser(u model.UserPayload) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var3 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var3 == nil {
			templ_7745c5c3_Var3 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templ.FromGoHTML(parsedCurrentUserTemplate, u).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}

const friendsTemplate = `
    <div id="friends-container" hx-get="/get-friends" hx-trigger="every 5s" hx-swap="outerHTML" class="my-4 mt-12">
        <div class="friends mr-4 ml-4">
            <h2 class="text-2xl font-bold mb-4 text-white" style="text-shadow: 2px 2px 2px #53a765;">Friends</h2>
            <ul class="space-y-4">
                {{range .}}
                <li class="rounded-lg shadow p-4 relative overflow-hidden mb-3">
                    <!-- Background Album Art (Blurred) -->
                    <div class="absolute inset-0 -z-10 shadow-md bg-cover bg-center blur-xl" style="background-image: url('{{.CurrentAlbumArt}}');"></div>
                    
                    <!-- Content Container -->
                    <div class="relative z-10 flex mb-2">
                        <!-- Profile Picture -->
                        <img src="{{.ProfilePicture}}" alt="Friend's Profile Picture" class="w-16 h-16 rounded-full mr-4"/>
                        <div>
                            <!-- Username -->
                            <h3 class="text-xl font-semibold text-white nunito-bold">{{.Username}}</h3>
                            <p class="text-sm text-white nunito-medium-italic">Currently Playing</p>
                        </div>
                    </div>

                    <!-- Album Art Beside Text -->
                    <div class="relative z-10 ml-20">
                        <img src="{{.CurrentAlbumArt}}" alt="Friend's Album Art" class="w-36 h-36 mb-2"/>
                        <div class="bg-black bg-opacity-20 backdrop-blur-lg rounded w-36">
                            <a href="{{.CurrentSongUrl}}" target="_blank" class="text-md text-white hover:underline nunito-bold-italic"><p class="text-center">{{.CurrentSongName}}</p></a>
                            <p class="text-sm text-white text-center nunito-medium">by {{.CurrentArtistName}}</p>
                            <p class="text-sm text-white text-center nunito-semibold">{{.CurrentAlbumName}}</p>
                        </div>
                    </div>
                </li>
                {{end}}
            </ul>
        </div>
    </div>
`

var parsedFriendsTemplate = template.Must(template.New("friends_template").Parse(friendsTemplate))

func Friends(friends []model.UserPayload) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var4 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var4 == nil {
			templ_7745c5c3_Var4 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templ.FromGoHTML(parsedFriendsTemplate, friends).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}

func NotificationBell() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var5 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var5 == nil {
			templ_7745c5c3_Var5 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<a href=\"/notifications\" class=\"relative items-center inline-flex\"><button id=\"notificationButton\" class=\"relative inline-flex items-center text-sm font-medium text-center text-gray-500 hover:text-gray-900 focus:outline-none dark:hover:text-white dark:text-gray-400\" type=\"button\"><svg class=\"w-5 h-5\" aria-hidden=\"true\" xmlns=\"http://www.w3.org/2000/svg\" fill=\"currentColor\" viewBox=\"0 0 14 20\"><path d=\"M12.133 10.632v-1.8A5.406 5.406 0 0 0 7.979 3.57.946.946 0 0 0 8 3.464V1.1a1 1 0 0 0-2 0v2.364a.946.946 0 0 0 .021.106 5.406 5.406 0 0 0-4.154 5.262v1.8C1.867 13.018 0 13.614 0 14.807 0 15.4 0 16 .538 16h12.924C14 16 14 15.4 14 14.807c0-1.193-1.867-1.789-1.867-4.175ZM3.823 17a3.453 3.453 0 0 0 6.354 0H3.823Z\"></path></svg><div class=\"absolute block w-3 h-3 bg-red-500 border-2 border-white rounded-full -top-0.5 start-2.5 dark:border-gray-900\"></div></button></a>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
