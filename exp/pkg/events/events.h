//go:build darwin

#ifndef _events_h
#define _events_h

extern void applicationEventHandler(unsigned int);
extern void windowEventHandler(unsigned int, unsigned int);

#define EventApplicationDidBecomeActive 0
#define EventApplicationDidChangeBackingProperties 1
#define EventApplicationDidChangeEffectiveAppearance 2
#define EventApplicationDidChangeIcon 3
#define EventApplicationDidChangeOcclusionState 4
#define EventApplicationDidChangeScreenParameters 5
#define EventApplicationDidChangeStatusBarFrame 6
#define EventApplicationDidChangeStatusBarOrientation 7
#define EventApplicationDidFinishLaunching 8
#define EventApplicationDidHide 9
#define EventApplicationDidResignActive 10
#define EventApplicationDidUnhide 11
#define EventApplicationDidUpdate 12
#define EventApplicationWillBecomeActive 13
#define EventApplicationWillFinishLaunching 14
#define EventApplicationWillHide 15
#define EventApplicationWillResignActive 16
#define EventApplicationWillTerminate 17
#define EventApplicationWillUnhide 18
#define EventApplicationWillUpdate 19
#define EventWindowDidBecomeKey 20
#define EventWindowDidBecomeMain 21
#define EventWindowDidBeginSheet 22
#define EventWindowDidChangeAlpha 23
#define EventWindowDidChangeBackingLocation 24
#define EventWindowDidChangeBackingProperties 25
#define EventWindowDidChangeCollectionBehavior 26
#define EventWindowDidChangeEffectiveAppearance 27
#define EventWindowDidChangeOcclusionState 28
#define EventWindowDidChangeOrderingMode 29
#define EventWindowDidChangeScreen 30
#define EventWindowDidChangeScreenParameters 31
#define EventWindowDidChangeScreenProfile 32
#define EventWindowDidChangeScreenSpace 33
#define EventWindowDidChangeScreenSpaceProperties 34
#define EventWindowDidChangeSharingType 35
#define EventWindowDidChangeSpace 36
#define EventWindowDidChangeSpaceOrderingMode 37
#define EventWindowDidChangeTitle 38
#define EventWindowDidChangeToolbar 39
#define EventWindowDidChangeVisibility 40
#define EventWindowDidClose 41
#define EventWindowDidDeminiaturize 42
#define EventWindowDidEndSheet 43
#define EventWindowDidEnterFullScreen 44
#define EventWindowDidEnterVersionBrowser 45
#define EventWindowDidExitFullScreen 46
#define EventWindowDidExitVersionBrowser 47
#define EventWindowDidExpose 48
#define EventWindowDidFocus 49
#define EventWindowDidMiniaturize 50
#define EventWindowDidMove 51
#define EventWindowDidOrderOffScreen 52
#define EventWindowDidOrderOnScreen 53
#define EventWindowDidResignKey 54
#define EventWindowDidResignMain 55
#define EventWindowDidResize 56
#define EventWindowDidUnfocus 57
#define EventWindowDidUpdate 58
#define EventWindowDidUpdateAlpha 59
#define EventWindowDidUpdateCollectionBehavior 60
#define EventWindowDidUpdateCollectionProperties 61
#define EventWindowDidUpdateShadow 62
#define EventWindowDidUpdateTitle 63
#define EventWindowDidUpdateToolbar 64
#define EventWindowDidUpdateVisibility 65
#define EventWindowWillBecomeKey 66
#define EventWindowWillBecomeMain 67
#define EventWindowWillBeginSheet 68
#define EventWindowWillChangeOrderingMode 69
#define EventWindowWillClose 70
#define EventWindowWillDeminiaturize 71
#define EventWindowWillEnterFullScreen 72
#define EventWindowWillEnterVersionBrowser 73
#define EventWindowWillExitFullScreen 74
#define EventWindowWillExitVersionBrowser 75
#define EventWindowWillFocus 76
#define EventWindowWillMiniaturize 77
#define EventWindowWillMove 78
#define EventWindowWillOrderOffScreen 79
#define EventWindowWillOrderOnScreen 80
#define EventWindowWillResignMain 81
#define EventWindowWillResize 82
#define EventWindowWillUnfocus 83
#define EventWindowWillUpdate 84
#define EventWindowWillUpdateAlpha 85
#define EventWindowWillUpdateCollectionBehavior 86
#define EventWindowWillUpdateCollectionProperties 87
#define EventWindowWillUpdateShadow 88
#define EventWindowWillUpdateTitle 89
#define EventWindowWillUpdateToolbar 90
#define EventWindowWillUpdateVisibility 91
#define EventWindowWillUseStandardFrame 92


#endif