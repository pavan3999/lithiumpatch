// # Hyphenation
//
// Add an option to enable/disable hyphenation.
package patches

import . "github.com/pgaskin/lithiumpatch/patches/patchdef"

func init() {
	Register("hyphenation",

		/* ---------------- JS (safe, unchanged) ---------------- */

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
				"\n"+`        style += hyphenation ? '-webkit-hyphens: auto; hyphens: auto;' : '-webkit-hyphens: none; hyphens: none;';`,
			),
		),

		/* ---------------- BookView / ContentView stubs ---------------- */

		PatchFiles(
	[]string{
		"smali/com/faultexception/reader/content/BookView.smali",
		"smali/com/faultexception/reader/content/ContentView.smali",
	},
	ReplaceStringPrepend(
	`.end class`,
	FixIndent("\n"+`
		.method public setHyphenation(Z)V
			.locals 0
			return-void
		.end method
	`),
   ),
),

		/* ---------------- HtmlContentView forwarding ---------------- */

		PatchFile("smali/com/faultexception/reader/content/HtmlContentView.smali",
			ReplaceStringPrepend(
				FixIndent("\n"+`
				.method public setTextAlign(I)V
				`),
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

		/* ---------------- EPubBookView state ---------------- */

		PatchFile("smali/com/faultexception/reader/content/EPubBookView.smali",
			ReplaceStringPrepend(
				"\n"+`.field private mTextAlign:I`,
				"\n"+`.field private mHyphenation:Z`,
			),
			InMethod("setTextAlign(I)V",
				ReplaceStringAppend(
					"\n",
					FixIndent("\n"+`
						iget-object v0, p0, Lcom/faultexception/reader/content/EPubBookView;->mContentView:Lcom/faultexception/reader/content/ContentView;
						if-eqz v0, :end
						invoke-virtual {v0, p1}, Lcom/faultexception/reader/content/ContentView;->setTextAlign(I)V
					:end
					`),
				),
			),
			InMethod("setHyphenation(Z)V",
				ReplaceStringAppend(
					"\n",
					FixIndent("\n"+`
						iput-boolean p1, p0, Lcom/faultexception/reader/content/EPubBookView;->mHyphenation:Z
						iget-object v0, p0, Lcom/faultexception/reader/content/EPubBookView;->mContentView:Lcom/faultexception/reader/content/ContentView;
						if-eqz v0, :end
						invoke-virtual {v0, p1}, Lcom/faultexception/reader/content/ContentView;->setHyphenation(Z)V
					:end
					`),
				),
			),
		),

		/* ---------------- HtmlContentWebView (ENTRY injection only) ---------------- */

		PatchFile("smali/com/faultexception/reader/content/HtmlContentWebView.smali",
			ReplaceStringPrepend(
				"\n"+`.field private mTextAlign:I`,
				"\n"+`.field private mHyphenation:Z`,
			),
			InMethod("setHyphenation(Z)V",
				ReplaceStringAppend(
					"\n",
					FixIndent("\n"+`
						iput-boolean p1, p0, Lcom/faultexception/reader/content/HtmlContentWebView;->mHyphenation:Z
						invoke-virtual {p0}, Lcom/faultexception/reader/content/HtmlContentWebView;->injectDisplaySettingsIfReady()V
					`),
				),
			),
		),

		/* ---------------- ReaderActivity preference hookup ---------------- */

		PatchFile("smali/com/faultexception/reader/ReaderActivity.smali",
			InMethod("updateFeaturesForBookView()V",
				ReplaceStringAppend(
					"\n",
					FixIndent("\n"+`
						iget-object v0, p0, Lcom/faultexception/reader/ReaderActivity;->mPrefs:Landroid/content/SharedPreferences;
						const-string v1, "hyphenation"
						const/4 v2, 0x1
						invoke-interface {v0, v1, v2}, Landroid/content/SharedPreferences;->getBoolean(Ljava/lang/String;Z)Z
						move-result v0
						iget-object v1, p0, Lcom/faultexception/reader/ReaderActivity;->mBookView:Lcom/faultexception/reader/content/BookView;
						invoke-virtual {v1, v0}, Lcom/faultexception/reader/content/BookView;->setHyphenation(Z)V
					`),
				),
			),
		),

		/* ---------------- Preferences UI ---------------- */

		PatchFile("res/xml/preferences.xml",
			ReplaceStringAppend(
				"\n"+`    <PreferenceCategory android:title="@string/pref_category_advanced">`,
				"\n"+`        <SwitchPreferenceCompat android:title="Use hyphenation" android:key="hyphenation" android:defaultValue="true" />`,
			),
		),
	)
}
