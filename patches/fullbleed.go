// # Full-bleed reader background
//
// Optionally expand the full-screen reader background into the status bar
// display cutout area.
package patches

import . "github.com/pgaskin/lithiumpatch/patches/patchdef"

func init() {
	Register("fullbleed",

		// Preference toggle
		PatchFile("res/xml/preferences.xml",
			ReplaceStringAppend(
				"\n"+`    <SwitchPreferenceCompat android:title="@string/pref_fullscreen_title" android:key="fullscreen" android:defaultValue="true" />`,
				"\n"+`    <SwitchPreferenceCompat android:title="Fullscreen reading full-bleed" android:key="fullscreen_bleed" android:defaultValue="true" />`,
			),
		),

		// ReaderActivity hooks
		PatchFile("smali/com/faultexception/reader/ReaderActivity.smali",

			// Inject helper method BEFORE setTheme
			ReplaceStringPrepend(
				FixIndent("\n"+`
				.method public setTheme(Lcom/faultexception/reader/themes/Theme;)V
				`),
				FixIndent("\n"+`
				.method private maybeSetDisplayCutoutBackgroundFromTheme(Lcom/faultexception/reader/themes/Theme;)V
					.locals 3

					# only if fullscreen enabled
					iget-boolean v0, p0, Lcom/faultexception/reader/ReaderActivity;->mFullscreenEnabled:Z
					if-eqz v0, :end

					# only if fullscreen full-bleed enabled
					invoke-static {p0}, Landroid/preference/PreferenceManager;->getDefaultSharedPreferences(Landroid/content/Context;)Landroid/content/SharedPreferences;
					move-result-object v0
					const-string v1, "fullscreen_bleed"
					const/4 v2, 0x0
					invoke-interface {v0, v1, v2}, Landroid/content/SharedPreferences;->getBoolean(Ljava/lang/String;Z)Z
					move-result v0
					if-eqz v0, :end

					# get theme background color
					iget v0, p1, Lcom/faultexception/reader/themes/Theme;->backgroundColor:I
					invoke-direct {p0, v0}, Lcom/faultexception/reader/ReaderActivity;->maybeSetDisplayCutoutBackground(I)V

				:end
					return-void
				.end method
				`),
			),

			// Call helper at END of setTheme (safe anchor)
			InMethod("setTheme(Lcom/faultexception/reader/themes/Theme;)V",
	ReplaceStringAppend(
		"\n",
		FixIndent("\n"+`
			invoke-direct {p0, p1}, Lcom/faultexception/reader/ReaderActivity;->maybeSetDisplayCutoutBackgroundFromTheme(Lcom/faultexception/reader/themes/Theme;)V
		`),
	),
),
		),

		// DisplayCutoutFrameLayout helper
		PatchFile("smali/com/faultexception/reader/widget/DisplayCutoutFrameLayout.smali",
			ReplaceStringPrepend(
				FixIndent("\n"+`
				.method public setInsetCutout(I)V
				`),
				FixIndent("\n"+`
				.method public setInsetCutoutColor(I)V
					.locals 1

					iget-boolean v0, p0, Lcom/faultexception/reader/widget/DisplayCutoutFrameLayout;->mPaintCutout:Z
					if-eqz v0, :end

					new-instance v0, Landroid/graphics/drawable/ColorDrawable;
					invoke-direct {v0, p1}, Landroid/graphics/drawable/ColorDrawable;-><init>(I)V
					iput-object v0, p0, Lcom/faultexception/reader/widget/DisplayCutoutFrameLayout;->mColor:Landroid/graphics/drawable/ColorDrawable;

				:end
					return-void
				.end method
				`),
			),
		),
	)
}
