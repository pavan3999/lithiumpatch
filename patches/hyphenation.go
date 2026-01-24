// # Hyphenation
//
// Add an option to enable/disable hyphenation.
package patches

import . "github.com/pgaskin/lithiumpatch/patches/patchdef"

func init() {
	Register("hyphenation",

		// --------------------
		// JS: epub.js
		// --------------------
		PatchFile("assets/js/epub.js",
			ReplaceStringPrepend(
				"\n"+`    var textAlign = void 0;`,
				"\n"+`    var hyphenation = void 0;`,
			),
			ReplaceStringPrepend(
				"\n"+`        setTextAlign: setTextAlign`,
				"\n"+`        setHyphenation: setHyphenation,`,
			),
			ReplaceStringPrepend(
				FixIndent("\n"+`
					function setTextAlign(align) {
						textAlign = align;
						updateStyleElement();
						reflowIfNecessary();
					}
				`),
				FixIndent("\n"+`
					function setHyphenation(hyp) {
						hyphenation = hyp;
						updateStyleElement();
						reflowIfNecessary();
					}
				`),
			),
			ReplaceStringPrepend(
				"\n"+`        styleElement.innerText = specificitySelector + ' * { ' + style + ' }';`,
				"\n"+`        style += hyphenation ? '-webkit-hyphens:auto;hyphens:auto;' : '-webkit-hyphens:none;hyphens:none;';`,
			),
		),

		// --------------------
		// BookView / ContentView
		// (interfaces only)
		// --------------------
		PatchFiles(
			[]string{
				"smali/com/faultexception/reader/content/BookView.smali",
				"smali/com/faultexception/reader/content/ContentView.smali",
			},
			ReplaceStringAppend(
				"",
				FixIndent("\n"+`
					.method public setHyphenation(Z)V
						.locals 0
						return-void
					.end method
				`),
			),
		),

		// --------------------
		// HtmlContentView
		// --------------------
		PatchFile("smali/com/faultexception/reader/content/HtmlContentView.smali",
			ReplaceStringAppend(
				"",
				FixIndent("\n"+`
					.method public setHyphenation(Z)V
						.locals 1
						iget-object v0, p0, Lcom/faultexception/reader/content/HtmlContentView;->mContentWebView:Lcom/faultexception/reader/content/HtmlContentWebView;
						invoke-virtual {v0, p1}, Lcom/faultexception/reader/content/HtmlContentWebView;->setHyphenation(Z)V
						return-void
					.end method
				`),
			),
		),

		// --------------------
		// EPubBookView (METHOD DID NOT EXIST → ADD IT)
		// --------------------
		PatchFile("smali/com/faultexception/reader/content/EPubBookView.smali",
			ReplaceStringPrepend(
				"\n"+`.field private mTextAlign:I`,
				"\n"+`.field private mHyphenation:Z`,
			),
			ReplaceStringAppend(
				"",
				FixIndent("\n"+`
					.method public setHyphenation(Z)V
						.locals 1

						iput-boolean p1, p0, Lcom/faultexception/reader/content/EPubBookView;->mHyphenation:Z

						iget-object v0, p0, Lcom/faultexception/reader/content/EPubBookView;->mContentView:Lcom/faultexception/reader/content/ContentView;
						if-eqz v0, :end

						invoke-virtual {v0, p1}, Lcom/faultexception/reader/content/ContentView;->setHyphenation(Z)V

					:end
						return-void
					.end method
				`),
			),
		),

		// --------------------
		// HtmlContentWebView
		// --------------------
		PatchFile("smali/com/faultexception/reader/content/HtmlContentWebView.smali",
			ReplaceStringPrepend(
				"\n"+`.field private mTextAlign:I`,
				"\n"+`.field private mHyphenation:Z`,
			),
			ReplaceStringAppend(
				"",
				FixIndent("\n"+`
					.method public setHyphenation(Z)V
						.locals 2

						iput-boolean p1, p0, Lcom/faultexception/reader/content/HtmlContentWebView;->mHyphenation:Z

						iget-object v0, p0, Lcom/faultexception/reader/content/HtmlContentWebView;->mUrl:Ljava/lang/String;
						if-eqz v0, :end

						iget-boolean v0, p0, Lcom/faultexception/reader/content/HtmlContentWebView;->mDisplaySettingsInjected:Z
						if-eqz v0, :end

						new-instance v0, Ljava/lang/StringBuilder;
						invoke-direct {v0}, Ljava/lang/StringBuilder;-><init>()V
						const-string v1, "LithiumJs.setHyphenation("
						invoke-virtual {v0, v1}, Ljava/lang/StringBuilder;->append(Ljava/lang/String;)Ljava/lang/StringBuilder;
						invoke-virtual {v0, p1}, Ljava/lang/StringBuilder;->append(Z)Ljava/lang/StringBuilder;
						const-string p1, ")"
						invoke-virtual {v0, p1}, Ljava/lang/StringBuilder;->append(Ljava/lang/String;)Ljava/lang/StringBuilder;
						invoke-virtual {v0}, Ljava/lang/StringBuilder;->toString()Ljava/lang/String;
						move-result-object p1
						invoke-virtual {p0, p1}, Lcom/faultexception/reader/content/HtmlContentWebView;->executeJavascript(Ljava/lang/String;)V

					:end
						return-void
					.end method
				`),
			),
		),

		// --------------------
		// ReaderActivity → apply preference
		// --------------------
		PatchFile("smali/com/faultexception/reader/ReaderActivity.smali",
			InMethod("updateFeaturesForBookView()V",
				ReplaceStringAppend(
					"\n",
					FixIndent("\n"+`
						iget-object v0, p0, Lcom/faultexception/reader/ReaderActivity;->mPrefs:Landroid/content/SharedPreferences;
						const/4 v1, 0x1
						const-string v2, "hyphenation"
						invoke-interface {v0, v2, v1}, Landroid/content/SharedPreferences;->getBoolean(Ljava/lang/String;Z)Z
						move-result v0

						iget-object v1, p0, Lcom/faultexception/reader/ReaderActivity;->mBookView:Lcom/faultexception/reader/content/BookView;
						invoke-virtual {v1, v0}, Lcom/faultexception/reader/content/BookView;->setHyphenation(Z)V
					`),
				),
			),
		),

		// --------------------
		// Preferences
		// --------------------
		PatchFile("res/xml/preferences.xml",
			ReplaceStringAppend(
				"\n"+`    <PreferenceCategory android:title="@string/pref_category_advanced">`,
				"\n"+`        <SwitchPreferenceCompat android:title="Use hyphenation" android:key="hyphenation" android:defaultValue="true" />`,
			),
		),
	)
}
